// Catalog affiche les produits. Il n'appelle personne. L'autorisation se fait
// à deux niveaux : handshake mTLS (membre du trust domain) puis politique par
// route, l'identité de l'appelant étant lue dans son certificat.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const trustDomain = "example.org"

const (
	gatewayID = "spiffe://example.org/ns/shop/sa/gateway"
)

// policy : qui peut appeler quelle route.
//
//	/products : lister les produits — réservé à la Gateway.
//	/_calls   : consulter le journal — réservé à la Gateway (pour le front).
var policy = map[string][]string{
	"/products": {gatewayID},
	"/_calls":   {gatewayID},
}

var products = []map[string]any{
	{"id": 1, "name": "Clé USB chiffrée", "price": 29},
	{"id": 2, "name": "YubiKey 5", "price": 55},
	{"id": 3, "name": "Câble Ethernet blindé", "price": 12},
}

type callEntry struct {
	Caller    string    `json:"caller"`
	Route     string    `json:"route"`
	Allowed   bool      `json:"allowed"`
	Timestamp time.Time `json:"timestamp"`
}

type callLog struct {
	mu      sync.Mutex
	entries []callEntry
}

func (l *callLog) record(e callEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, e)
	if len(l.entries) > 50 {
		l.entries = l.entries[len(l.entries)-50:]
	}
}

func (l *callLog) snapshot() []callEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]callEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

func main() {
	ctx := context.Background()
	socketPath := getenv("SPIFFE_ENDPOINT_SOCKET", "unix:///run/spire/sockets/spire-agent.sock")

	source, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Fatalf("catalog: impossible d'obtenir un SVID via le Workload API: %v", err)
	}
	defer source.Close()

	svid, err := source.GetX509SVID()
	if err != nil {
		log.Fatalf("catalog: pas de SVID disponible: %v", err)
	}
	log.Printf("catalog: identité %s", svid.ID)

	logbook := &callLog{}

	go serveMetrics()

	mux := http.NewServeMux()
	mux.HandleFunc("/products", guard("/products", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		log.Printf("catalog: produits servis à %s", caller)
		writeJSON(w, http.StatusOK, products)
	}))
	mux.HandleFunc("/_calls", guard("/_calls", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		writeJSON(w, http.StatusOK, logbook.snapshot())
	}))

	td, _ := spiffeid.TrustDomainFromString(trustDomain)
	tlsConfig := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeMemberOf(td))

	server := &http.Server{Addr: ":8443", Handler: mux, TLSConfig: tlsConfig}
	log.Println("catalog: écoute en mTLS sur :8443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func guard(route string, logbook *callLog, h func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller := callerID(r)
		allowed := slices.Contains(policy[route], caller)
		recordAuthz(allowed)
		// /_calls est une route d'introspection (lecture du journal par le front) ;
		// on ne la journalise pas pour ne pas noyer les vrais appels métier.
		if route != "/_calls" {
			logbook.record(callEntry{Caller: caller, Route: route, Allowed: allowed, Timestamp: time.Now()})
		}
		if !allowed {
			log.Printf("catalog: %s REFUSÉ sur %s", caller, route)
			recordRequest(route, http.StatusForbidden)
			writeJSON(w, http.StatusForbidden, map[string]any{
				"allowed": false, "route": route, "caller": caller,
				"reason": "identité non autorisée pour cette route",
			})
			return
		}
		recordRequest(route, http.StatusOK)
		h(w, r, caller)
	}
}

func callerID(r *http.Request) string {
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 || len(r.TLS.PeerCertificates[0].URIs) == 0 {
		return "(inconnu)"
	}
	id, err := spiffeid.FromURI(r.TLS.PeerCertificates[0].URIs[0])
	if err != nil {
		return "(invalide)"
	}
	return id.String()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
