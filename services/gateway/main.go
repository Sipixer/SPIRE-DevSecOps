// Gateway : porte d'entrée HTTP du navigateur, relaie vers les services (mTLS Istio).
package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed static
var staticFiles embed.FS

func main() {
	ordersURL := getenv("ORDERS_URL", "http://orders.shop.svc.cluster.local:8080")
	catalogURL := getenv("CATALOG_URL", "http://catalog.shop.svc.cluster.local:8080")
	paymentsURL := getenv("PAYMENTS_URL", "http://payments.shop.svc.cluster.local:8080")
	analyticsURL := getenv("ANALYTICS_URL", "http://analytics.shop.svc.cluster.local:8080")

	client := &http.Client{Timeout: 5 * time.Second}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, ordersURL+"/order", "POST")
	})

	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, catalogURL+"/products", "GET")
	})

	// /pay est réservé à Orders : l'AuthorizationPolicy de Payments refuse la Gateway (403).
	mux.HandleFunc("/api/forbidden", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, paymentsURL+"/pay", "POST")
	})

	mux.HandleFunc("/api/event", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, analyticsURL+"/event", "POST")
	})

	mux.HandleFunc("/api/calls", func(w http.ResponseWriter, r *http.Request) {
		out := map[string]any{
			"orders":   fetchJSON(client, ordersURL+"/_calls"),
			"catalog":  fetchJSON(client, catalogURL+"/_calls"),
			"payments": fetchJSON(client, paymentsURL+"/_calls"),
		}
		writeJSON(w, http.StatusOK, out)
	})

	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("gateway: impossible de monter le front: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticRoot)))

	addr := getenv("LISTEN_ADDR", ":8080")
	log.Printf("gateway: écoute en HTTP sur %s (front + API navigateur)", addr)
	log.Fatal(http.ListenAndServe(addr, mux)) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
}

// proxyJSON recopie la réponse au navigateur en conservant le statut (un 403 arrive tel quel).
func proxyJSON(w http.ResponseWriter, client *http.Client, url, method string) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"allowed": false, "error": err.Error()})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"allowed": false, "error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func fetchJSON(client *http.Client, url string) any {
	resp, err := client.Get(url)
	if err != nil {
		return map[string]string{"error": err.Error()}
	}
	defer resp.Body.Close()
	var v any
	json.NewDecoder(resp.Body).Decode(&v)
	return v
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
