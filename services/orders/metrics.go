// Observabilité : compteurs Prometheus tenus à la main (pas de dépendance
// externe) et exposés en HTTP clair sur :9100 via GET /metrics.
package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

const serviceName = "orders"

// counter est un ensemble de compteurs étiquetés, protégé par un mutex.
// La clé est la ligne de labels déjà formatée (ex: path="/pay",code="200").
type counter struct {
	mu     sync.Mutex
	name   string
	help   string
	values map[string]int
}

func newCounter(name, help string) *counter {
	return &counter{name: name, help: help, values: map[string]int{}}
}

func (c *counter) inc(labels string) {
	c.mu.Lock()
	c.values[labels]++
	c.mu.Unlock()
}

// write rend le compteur au format texte Prometheus.
func (c *counter) write(w http.ResponseWriter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Fprintf(w, "# HELP %s %s\n# TYPE %s counter\n", c.name, c.help, c.name)
	keys := make([]string, 0, len(c.values))
	for k := range c.values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s{%s} %d\n", c.name, k, c.values[k])
	}
}

var (
	httpRequests   = newCounter("http_requests_total", "Requêtes HTTP traitées.")
	authzDecisions = newCounter("authz_decisions_total", "Décisions d'autorisation par route.")
)

func recordRequest(path string, code int) {
	httpRequests.inc(fmt.Sprintf("service=%q,path=%q,code=%q", serviceName, path, strconv.Itoa(code)))
}

func recordAuthz(allowed bool) {
	decision := "denied"
	if allowed {
		decision = "allowed"
	}
	authzDecisions.inc(fmt.Sprintf("service=%q,decision=%q", serviceName, decision))
}

// serveMetrics lance le listener d'observabilité en clair sur :9100.
func serveMetrics() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		httpRequests.write(w)
		authzDecisions.write(w)
	})
	log.Println(serviceName + ": métriques exposées en HTTP sur :9100 (/metrics)")
	log.Fatal(http.ListenAndServe(":9100", mux)) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
}
