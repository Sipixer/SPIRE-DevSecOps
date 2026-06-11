// Orders gère les commandes. Il est serveur mTLS (reçoit la Gateway) et client
// mTLS (appelle Payments). L'autorisation se fait au handshake (membre du trust
// domain) puis par route, l'identité étant lue dans le certificat client.
package main

import (
	"context"
	"encoding/json"
	"io"
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
	gatewayID  = "spiffe://example.org/ns/shop/sa/gateway"
	paymentsID = "spiffe://example.org/ns/shop/sa/payments"
)

// policy : qui peut appeler quelle route d'Orders.
//   /order  : passer une commande — réservé à la Gateway.
//   /_calls : consulter le journal — réservé à la Gateway (pour le front).
var policy = map[string][]string{
	"/order":  {gatewayID},
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
	paymentsURL := getenv("PAYMENTS_URL", "https://payments.shop.svc.cluster.local:8443/pay")

	source, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Fatalf("orders: impossible d'obtenir un SVID via le Workload API: %v", err)
	}
	defer source.Close()

	svid, err := source.GetX509SVID()
	if err != nil {
		log.Fatalf("orders: pas de SVID disponible: %v", err)
	}
	log.Printf("orders: identité %s", svid.ID)

	// Client mTLS vers Payments : présente le SVID d'Orders et n'accepte de
	// parler qu'au serveur dont l'identité est exactement celle de Payments.
	paymentsSPIFFE, err := spiffeid.FromString(paymentsID)
	if err != nil {
		log.Fatalf("orders: SPIFFE ID Payments invalide: %v", err)
	}
	clientTLS := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeID(paymentsSPIFFE))
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: clientTLS},
		Timeout:   5 * time.Second,
	}

	logbook := &callLog{}

	mux := http.NewServeMux()
	mux.HandleFunc("/order", guard("/order", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		resp, err := client.Post(paymentsURL, "application/json", nil)
		if err != nil {
			log.Printf("orders: appel Payments échoué: %v", err)
			writeJSON(w, http.StatusBadGateway, map[string]any{
				"order": "échec", "caller": caller, "error": err.Error(), "allowed": false,
			})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		log.Printf("orders: commande de %s, Payments a répondu %d", caller, resp.StatusCode)
		writeJSON(w, http.StatusOK, map[string]any{
			"order": "ok", "caller": caller,
			"payments_status": resp.StatusCode, "payments_body": json.RawMessage(body),
		})
	}))
	mux.HandleFunc("/_calls", guard("/_calls", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		writeJSON(w, http.StatusOK, logbook.snapshot())
	}))

	td, _ := spiffeid.TrustDomainFromString(trustDomain)
	serverTLS := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeMemberOf(td))
	server := &http.Server{Addr: ":8443", Handler: mux, TLSConfig: serverTLS}
	log.Println("orders: écoute en mTLS sur :8443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

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
			log.Printf("orders: %s REFUSÉ sur %s", caller, route)
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
