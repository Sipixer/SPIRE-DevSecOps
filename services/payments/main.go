// Payments est un service sensible. L'autorisation se fait à deux niveaux :
//   1. au handshake mTLS : on n'accepte que les identités du trust domain ;
//   2. par route, dans le code : chaque route déclare qui a le droit de l'appeler,
//      l'identité étant lue dans le certificat client (non falsifiable).
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

// Identités des services, pour rendre la politique lisible.
const (
	gatewayID  = "spiffe://example.org/ns/shop/sa/gateway"
	ordersID   = "spiffe://example.org/ns/shop/sa/orders"
	catalogID  = "spiffe://example.org/ns/shop/sa/catalog"
	paymentsID = "spiffe://example.org/ns/shop/sa/payments"
)

// policy déclare, pour chaque route, la liste blanche des appelants autorisés.
//   /pay     : déclencher un paiement — réservé à Orders.
//   /_calls  : consulter le journal d'appels — réservé à la Gateway (pour le front).
var policy = map[string][]string{
	"/pay":    {ordersID},
	"/_calls": {gatewayID},
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
		log.Fatalf("payments: impossible d'obtenir un SVID via le Workload API: %v", err)
	}
	defer source.Close()

	svid, err := source.GetX509SVID()
	if err != nil {
		log.Fatalf("payments: pas de SVID disponible: %v", err)
	}
	log.Printf("payments: identité %s", svid.ID)

	logbook := &callLog{}

	mux := http.NewServeMux()
	mux.HandleFunc("/pay", guard("/pay", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		log.Printf("payments: paiement autorisé pour %s", caller)
		writeJSON(w, http.StatusOK, map[string]string{"status": "paid", "caller": caller})
	}))
	mux.HandleFunc("/_calls", guard("/_calls", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		writeJSON(w, http.StatusOK, logbook.snapshot())
	}))

	// Au handshake, on accepte toute identité du trust domain ; le filtrage
	// fin par route est fait ensuite par guard().
	td, _ := spiffeid.TrustDomainFromString(trustDomain)
	tlsConfig := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeMemberOf(td))

	server := &http.Server{Addr: ":8443", Handler: mux, TLSConfig: tlsConfig}
	log.Println("payments: écoute en mTLS sur :8443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

// guard applique la politique de la route : il lit l'identité de l'appelant
// dans le certificat, vérifie la liste blanche, journalise, puis délègue.
func guard(route string, logbook *callLog, h func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller := callerID(r)
		allowed := slices.Contains(policy[route], caller)
		// /_calls est une route d'introspection (lecture du journal par le front) ;
		// on ne la journalise pas pour ne pas noyer les vrais appels métier.
		if route != "/_calls" {
			logbook.record(callEntry{Caller: caller, Route: route, Allowed: allowed, Timestamp: time.Now()})
		}
		if !allowed {
			log.Printf("payments: %s REFUSÉ sur %s", caller, route)
			writeJSON(w, http.StatusForbidden, map[string]any{
				"allowed": false, "route": route, "caller": caller,
				"reason": "identité non autorisée pour cette route",
			})
			return
		}
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
