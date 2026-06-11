// Gateway est la porte d'entrée. Le navigateur lui parle en HTTP simple ;
// elle relaie vers Orders et Catalog en mTLS, en présentant son SVID. Elle
// sert aussi le front et agrège les journaux d'appels de chaque service.
package main

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

//go:embed static
var staticFiles embed.FS

func main() {
	ctx := context.Background()
	socketPath := getenv("SPIFFE_ENDPOINT_SOCKET", "unix:///run/spire/sockets/spire-agent.sock")
	ordersURL := getenv("ORDERS_URL", "https://orders.shop.svc.cluster.local:8443")
	catalogURL := getenv("CATALOG_URL", "https://catalog.shop.svc.cluster.local:8443")
	paymentsURL := getenv("PAYMENTS_URL", "https://payments.shop.svc.cluster.local:8443")
	analyticsURL := getenv("ANALYTICS_URL", "https://analytics.shop.svc.cluster.local:8443")

	source, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Fatalf("gateway: impossible d'obtenir un SVID via le Workload API: %v", err)
	}
	defer source.Close()

	svid, err := source.GetX509SVID()
	if err != nil {
		log.Fatalf("gateway: pas de SVID disponible: %v", err)
	}
	log.Printf("gateway: identité %s", svid.ID)

	orders := mtlsClient(source, "spiffe://example.org/ns/shop/sa/orders")
	catalog := mtlsClient(source, "spiffe://example.org/ns/shop/sa/catalog")
	payments := mtlsClient(source, "spiffe://example.org/ns/shop/sa/payments")
	analytics := mtlsClient(source, "spiffe://example.org/ns/shop/sa/analytics")

	go serveMetrics()

	mux := http.NewServeMux()

	// --- Routes côté navigateur (HTTP simple) ---

	// Déclenche Orders -> Payments en mTLS.
	mux.HandleFunc("/api/order", instrument("/api/order", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, orders, ordersURL+"/order", "POST")
	}))

	// Récupère le catalogue auprès de Catalog en mTLS.
	mux.HandleFunc("/api/products", instrument("/api/products", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, catalog, catalogURL+"/products", "GET")
	}))

	// Démonstration du refus : la Gateway tente d'appeler /pay de Payments,
	// route réservée à Orders. Le handshake mTLS réussit (Gateway est dans le
	// trust domain) mais la politique par route renvoie 403. C'est le cœur du
	// zero-trust : être authentifié ne suffit pas, il faut être autorisé.
	mux.HandleFunc("/api/forbidden", instrument("/api/forbidden", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, payments, paymentsURL+"/pay", "POST")
	}))

	// Envoie un événement à Analytics en mTLS (identité attendue : analytics).
	mux.HandleFunc("/api/event", instrument("/api/event", func(w http.ResponseWriter, r *http.Request) {
		proxyJSON(w, analytics, analyticsURL+"/event", "POST")
	}))

	// Agrège les journaux d'appels des services pour le front. La Gateway est
	// la seule autorisée à lire /_calls sur chaque service.
	mux.HandleFunc("/api/calls", instrument("/api/calls", func(w http.ResponseWriter, r *http.Request) {
		out := map[string]any{
			"orders":   fetchJSON(orders, ordersURL+"/_calls"),
			"catalog":  fetchJSON(catalog, catalogURL+"/_calls"),
			"payments": fetchJSON(payments, paymentsURL+"/_calls"),
		}
		writeJSON(w, http.StatusOK, out)
	}))

	// Mode démo : génère un trafic cohérent (commandes OK, catalogue,
	// événements analytics, et le scénario de refus) pour peupler les
	// journaux et les métriques. Retourne un résumé {scenarios, allowed, denied}.
	mux.HandleFunc("/api/demo", instrument("/api/demo", func(w http.ResponseWriter, r *http.Request) {
		runDemo(w, orders, catalog, payments, analytics,
			ordersURL, catalogURL, paymentsURL, analyticsURL)
	}))

	// Le front statique (contenu du dossier static/, servi à la racine).
	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("gateway: impossible de monter le front: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticRoot)))

	log.Println("gateway: écoute en HTTP sur :8080 (front + API navigateur)")
	// HTTP volontaire ici : c'est la seule surface côté navigateur. Le TLS
	// public est terminé par l'ingress (Let's Encrypt) ; tout le trafic
	// inter-services, lui, est en mTLS SPIFFE. Exception tracée et limitée
	// à cette ligne plutôt que désactivée globalement.
	log.Fatal(http.ListenAndServe(":8080", mux)) // nosemgrep: go.lang.security.audit.net.use-tls.use-tls
}

// mtlsClient construit un client HTTP mTLS qui n'accepte que le serveur
// dont l'identité SPIFFE correspond à peerID.
func mtlsClient(source *workloadapi.X509Source, peerID string) *http.Client {
	id, err := spiffeid.FromString(peerID)
	if err != nil {
		log.Fatalf("gateway: SPIFFE ID pair invalide %q: %v", peerID, err)
	}
	cfg := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeID(id))
	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: cfg},
		Timeout:   5 * time.Second,
	}
}

// proxyJSON appelle une URL en mTLS et recopie la réponse au navigateur. Le
// statut HTTP est conservé : un 403 (refus de la politique par route) ou un
// échec de handshake mTLS arrivent tels quels au front.
func proxyJSON(w http.ResponseWriter, client *http.Client, url, method string) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"allowed": false, "error": err.Error()})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		// Le handshake mTLS lui-même a échoué (identité non membre du domaine).
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

// runDemo joue un scénario d'appels pour produire un trafic réaliste : des
// commandes et consultations qui réussissent, des événements analytics, et le
// scénario de refus (Gateway -> /pay, réservé à Orders). Compte les décisions.
func runDemo(w http.ResponseWriter, orders, catalog, payments, analytics *http.Client,
	ordersURL, catalogURL, paymentsURL, analyticsURL string) {
	steps := []struct {
		client *http.Client
		url    string
		method string
	}{
		{orders, ordersURL + "/order", "POST"},
		{catalog, catalogURL + "/products", "GET"},
		{analytics, analyticsURL + "/event", "POST"},
		{orders, ordersURL + "/order", "POST"},
		{catalog, catalogURL + "/products", "GET"},
		{payments, paymentsURL + "/pay", "POST"}, // refus attendu (403)
		{analytics, analyticsURL + "/event", "POST"},
		{orders, ordersURL + "/order", "POST"},
	}

	allowed, denied := 0, 0
	for _, s := range steps {
		if callOK(s.client, s.url, s.method) {
			allowed++
		} else {
			denied++
		}
	}
	writeJSON(w, http.StatusOK, map[string]int{
		"scenarios": len(steps), "allowed": allowed, "denied": denied,
	})
}

// callOK exécute un appel mTLS et renvoie true si le pair a répondu 2xx.
func callOK(client *http.Client, url, method string) bool {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode < 300
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
