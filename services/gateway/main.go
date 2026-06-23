// Gateway est la porte d'entrée. Le navigateur lui parle en HTTP simple ;
// elle relaie vers Orders, Catalog, Payments et Analytics. En v2, le mTLS
// inter-services est porté par le mesh (Linkerd) : la Gateway fait du HTTP
// clair vers chaque service, et le proxy chiffre/authentifie de façon
// transparente. Plus de client mTLS à construire à la main. Elle sert aussi le
// front et agrège les journaux d'appels de chaque service.
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

	// Un seul client HTTP suffit : le mesh ajoute le mTLS et n'autorise que les
	// destinations permises par les AuthorizationPolicy. L'identité de la Gateway
	// est portée par son ServiceAccount, présentée automatiquement par le proxy.
	client := &http.Client{Timeout: 5 * time.Second}

	// Les métriques (RPS, latence, codes HTTP, mTLS) sont produites par le proxy
	// Linkerd : plus de listener /metrics ni de compteurs maison côté app.

	mux := http.NewServeMux()

	// --- Routes côté navigateur (HTTP simple) ---

	// Déclenche Orders -> Payments via le mesh.
	mux.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, ordersURL+"/order", "POST")
	})

	// Récupère le catalogue auprès de Catalog via le mesh.
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, catalogURL+"/products", "GET")
	})

	// Démonstration du refus : la Gateway tente d'appeler /pay de Payments,
	// route réservée à Orders. Le mesh authentifie la Gateway (mTLS OK) mais
	// l'AuthorizationPolicy de Payments refuse l'appel : le proxy renvoie 403
	// AVANT que la requête n'atteigne l'app. C'est le cœur du zero-trust : être
	// authentifié ne suffit pas, il faut être autorisé — et c'est le mesh, pas
	// le code applicatif, qui le décide.
	mux.HandleFunc("/api/forbidden", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, paymentsURL+"/pay", "POST")
	})

	// Envoie un événement à Analytics via le mesh.
	mux.HandleFunc("/api/event", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, client, analyticsURL+"/event", "POST")
	})

	// Agrège les journaux d'appels des services pour le front.
	mux.HandleFunc("/api/calls", func(w http.ResponseWriter, r *http.Request) {
		out := map[string]any{
			"orders":   fetchJSON(client, ordersURL+"/_calls"),
			"catalog":  fetchJSON(client, catalogURL+"/_calls"),
			"payments": fetchJSON(client, paymentsURL+"/_calls"),
		}
		writeJSON(w, http.StatusOK, out)
	})

	// Le front statique (contenu du dossier static/, servi à la racine).
	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("gateway: impossible de monter le front: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticRoot)))

	addr := getenv("LISTEN_ADDR", ":8080")
	log.Printf("gateway: écoute en HTTP sur %s (front + API navigateur)", addr)
	// HTTP côté navigateur : le TLS public est terminé par l'ingress. Le trafic
	// inter-services est en mTLS via le mesh.
	log.Fatal(http.ListenAndServe(addr, mux)) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
}

// proxyJSON appelle une URL via le mesh et recopie la réponse au navigateur. Le
// statut HTTP est conservé : un 403 (refus d'une AuthorizationPolicy, émis par
// le proxy Linkerd) arrive tel quel au front.
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
