// Analytics : service applicatif Node.js placé derrière Envoy.
// Envoy porte le mTLS SPIFFE, vérifie l'identité du pair et la transmet ici
// via un header. L'app écoute en clair sur localhost : seule Envoy l'atteint.
"use strict";

const http = require("http");

const HOST = "127.0.0.1";
const PORT = 9090;
const METRICS_PORT = 9100;

let eventsTotal = 0;
const eventsByCaller = new Map();

// Extrait le SPIFFE ID de l'appelant transmis par Envoy. Envoy peut poser un
// header dédié (x-spiffe-id) ou le XFCC standard (x-forwarded-client-cert),
// qui contient un champ URI=spiffe://... pour le certificat client.
function callerFromHeaders(headers) {
  const direct = headers["x-spiffe-id"];
  if (direct) return direct;
  const xfcc = headers["x-forwarded-client-cert"];
  if (xfcc) {
    const match = xfcc.match(/URI=([^;,]+)/i);
    if (match) return match[1];
  }
  return "unknown";
}

function metrics() {
  const lines = [
    "# HELP analytics_events_total Nombre total d'événements reçus.",
    "# TYPE analytics_events_total counter",
    `analytics_events_total ${eventsTotal}`,
    "# HELP analytics_events_by_caller_total Événements reçus par identité SPIFFE appelante.",
    "# TYPE analytics_events_by_caller_total counter",
  ];
  for (const [caller, count] of eventsByCaller) {
    lines.push(`analytics_events_by_caller_total{caller="${caller}"} ${count}`);
  }
  return lines.join("\n") + "\n";
}

const server = http.createServer((req, res) => {
  if (req.method === "POST" && req.url === "/event") {
    const caller = callerFromHeaders(req.headers);
    eventsTotal++;
    eventsByCaller.set(caller, (eventsByCaller.get(caller) || 0) + 1);
    res.writeHead(200, { "content-type": "application/json" });
    res.end(JSON.stringify({ ok: true, caller, events: eventsTotal }));
    return;
  }
  if (req.method === "GET" && req.url === "/metrics") {
    res.writeHead(200, { "content-type": "text/plain; version=0.0.4" });
    res.end(metrics());
    return;
  }
  if (req.method === "GET" && req.url === "/healthz") {
    res.writeHead(200, { "content-type": "text/plain" });
    res.end("ok\n");
    return;
  }
  res.writeHead(404, { "content-type": "text/plain" });
  res.end("not found\n");
});

server.listen(PORT, HOST, () => {
  console.log(`analytics: écoute sur http://${HOST}:${PORT} (derrière Envoy)`);
});

// Endpoint /metrics exposé au cluster (HTTP clair, comme les services Go) pour
// que Prometheus le scrape. Pas de donnée sensible, jamais exposé hors cluster.
http
  .createServer((req, res) => {
    if (req.method === "GET" && req.url === "/metrics") {
      res.writeHead(200, { "content-type": "text/plain; version=0.0.4" });
      res.end(metrics());
      return;
    }
    res.writeHead(404);
    res.end();
  })
  .listen(METRICS_PORT, "0.0.0.0");
