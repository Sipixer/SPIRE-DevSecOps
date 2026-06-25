// Analytics : compteur d'événements (mTLS et identité de l'appelant portés par le sidecar Istio).
"use strict";

const http = require("http");

const HOST = "0.0.0.0";
const PORT = Number(process.env.LISTEN_PORT || 8080);

let eventsTotal = 0;
const eventsByCaller = new Map();

function callerFromHeaders(headers) {
  const xfcc = headers["x-forwarded-client-cert"];
  if (xfcc) {
    const match = xfcc.match(/URI=([^;,]+)/i);
    if (match) return match[1];
  }
  return "unknown";
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
  if (req.method === "GET" && req.url === "/healthz") {
    res.writeHead(200, { "content-type": "text/plain" });
    res.end("ok\n");
    return;
  }
  res.writeHead(404, { "content-type": "text/plain" });
  res.end("not found\n");
});

server.listen(PORT, HOST, () => {
  console.log(`analytics: écoute sur http://${HOST}:${PORT} (mTLS assuré par le sidecar Istio)`);
});
