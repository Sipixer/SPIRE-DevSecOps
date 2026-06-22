// Orders gère les commandes. Il reçoit la Gateway et appelle Payments. En v2,
// le mTLS entre services est automatique (mesh Linkerd) : Orders parle HTTP
// clair, et le proxy chiffre/authentifie pour lui. L'autorisation par route
// (qui peut appeler /order, qui peut appeler /pay) est déclarée hors du code
// dans des AuthorizationPolicy Linkerd (voir k8s/workloads/authz.yaml).
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

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
	// HTTP clair : le mTLS vers Payments est ajouté de façon transparente par le
	// proxy Linkerd. Plus de client TLS à construire à la main.
	paymentsURL := getenv("PAYMENTS_URL", "http://payments.shop.svc.cluster.local:8080/pay")
	client := &http.Client{Timeout: 5 * time.Second}

	logbook := &callLog{}

	mux := http.NewServeMux()
	mux.HandleFunc("/order", logged("/order", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
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
	mux.HandleFunc("/_calls", logged("/_calls", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		writeJSON(w, http.StatusOK, logbook.snapshot())
	}))

	addr := getenv("LISTEN_ADDR", ":8080")
	log.Printf("orders: écoute en HTTP sur %s (mTLS assuré par le mesh)", addr)
	log.Fatal(http.ListenAndServe(addr, mux)) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
}

// logged journalise l'appel puis délègue. L'autorisation est déjà faite par le
// proxy Linkerd en amont : toute requête qui arrive ici a été autorisée.
func logged(route string, logbook *callLog, h func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller := callerID(r)
		if route != "/_calls" {
			logbook.record(callEntry{Caller: caller, Route: route, Allowed: true, Timestamp: time.Now()})
		}
		h(w, r, caller)
	}
}

// callerID lit l'identité de l'appelant dans l'en-tête posé par le proxy Linkerd.
func callerID(r *http.Request) string {
	id := r.Header.Get("l5d-client-id")
	if id == "" {
		return "(inconnu)"
	}
	parts := strings.Split(id, ".")
	if len(parts) >= 2 {
		return "spiffe://" + parts[1] + "/sa/" + parts[0]
	}
	return id
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
