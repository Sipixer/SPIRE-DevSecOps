// Payments est un service sensible. En v2, l'identité et l'autorisation sont
// portées par le service mesh (Linkerd) :
//  1. le mTLS est automatique entre proxies (l'app parle HTTP clair en local) ;
//  2. l'autorisation par route est déclarée hors du code, dans des
//     AuthorizationPolicy Linkerd (voir k8s/workloads/authz.yaml). Une requête
//     qui n'est pas autorisée n'atteint JAMAIS ce code : le proxy la rejette.
//
// L'app garde un journal d'appels pour le front : l'identité de l'appelant est
// désormais fournie par Linkerd via l'en-tête l5d-client-id (non falsifiable,
// posé par le proxy après vérification du mTLS).
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
	mux.HandleFunc("/pay", logged("/pay", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		log.Printf("payments: paiement pour %s", caller)
		writeJSON(w, http.StatusOK, map[string]string{"status": "paid", "caller": caller})
	}))
	mux.HandleFunc("/_calls", logged("/_calls", logbook, func(w http.ResponseWriter, r *http.Request, caller string) {
		writeJSON(w, http.StatusOK, logbook.snapshot())
	}))

	addr := getenv("LISTEN_ADDR", ":8080")
	log.Printf("payments: écoute en HTTP sur %s (mTLS assuré par le mesh)", addr)
	log.Fatal(http.ListenAndServe(addr, mux)) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
}

// logged journalise l'appel puis délègue. Il ne décide plus de l'autorisation :
// c'est le proxy Linkerd qui l'a déjà fait en amont (AuthorizationPolicy). Toute
// requête qui arrive ici a donc été autorisée — d'où Allowed: true. Les refus
// sont visibles dans Linkerd Viz (réponses 403 émises par le proxy), pas ici.
func logged(route string, logbook *callLog, h func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller := callerID(r)
		// /_calls est une route d'introspection (lecture du journal par le front) ;
		// on ne la journalise pas pour ne pas noyer les vrais appels métier.
		if route != "/_calls" {
			logbook.record(callEntry{Caller: caller, Route: route, Allowed: true, Timestamp: time.Now()})
		}
		h(w, r, caller)
	}
}

// callerID lit l'identité de l'appelant dans l'en-tête posé par le proxy
// Linkerd. Le format est "<sa>.<ns>.serviceaccount.identity.linkerd.cluster.local".
// On le réécrit en forme SPIFFE-like pour rester lisible côté front.
func callerID(r *http.Request) string {
	id := r.Header.Get("l5d-client-id")
	if id == "" {
		return "(inconnu)"
	}
	// gateway.shop.serviceaccount... -> spiffe://shop/sa/gateway (lisible).
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
