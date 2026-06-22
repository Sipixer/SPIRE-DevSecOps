// Catalog affiche les produits. Il n'appelle personne. En v2, le mTLS et
// l'autorisation sont portés par le mesh (Linkerd) : l'app parle HTTP clair,
// et la politique « qui peut appeler /products » est déclarée hors du code dans
// des AuthorizationPolicy Linkerd (voir k8s/workloads/authz.yaml). Une requête
// non autorisée est rejetée par le proxy avant d'atteindre ce code.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

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
	logbook := &callLog{}

	mux := http.NewServeMux()
	mux.HandleFunc("/products", logged("/products", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		log.Printf("catalog: produits servis à %s", caller)
		writeJSON(w, http.StatusOK, products)
	}))
	mux.HandleFunc("/_calls", logged("/_calls", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		writeJSON(w, http.StatusOK, logbook.snapshot())
	}))

	addr := getenv("LISTEN_ADDR", ":8080")
	log.Printf("catalog: écoute en HTTP sur %s (mTLS assuré par le mesh)", addr)
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
