// -------------------------------------------------------
// tools/serve_static.go
// -------------------------------------------------------
// Purpose Summary:
//   - Lightweight HTTP static file server for frontend files.
//   - Intended for local testing or airgapped deployment.
// Audit:
//   - Logs start and all HTTP requests with UTC timestamp.
//   - Fails fast on port conflict or missing frontend dir.
// -------------------------------------------------------

package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

const port = "7777"
const rootDir = "./frontend"

// -------------------------------------------------------
// func utcNow()
// -------------------------------------------------------
// Purpose:
//   - Provides consistent ISO UTC timestamps for audit logs.
// Audit:
//   - Used for all log statements in this static server.
// -------------------------------------------------------
func utcNow() string {
    return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

// -------------------------------------------------------
// func main()
// -------------------------------------------------------
// Purpose:
//   - Starts HTTP server to serve static frontend assets.
// Audit:
//   - Logs request paths and errors explicitly.
// -------------------------------------------------------
func main() {
    if _, err := os.Stat(rootDir); os.IsNotExist(err) {
        log.Fatalf("[ERROR] %s frontend folder not found: %s\n", utcNow(), rootDir)
    }

    http.Handle("/", logMiddleware(http.FileServer(http.Dir(rootDir))))

    log.Printf("[INFO] %s Serving %s on http://localhost:%s\n", utcNow(), rootDir, port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatalf("[ERROR] %s Server failed: %s\n", utcNow(), err.Error())
    }
}

// -------------------------------------------------------
// func logMiddleware(h http.Handler)
// -------------------------------------------------------
// Purpose:
//   - Logs all HTTP requests to stdout with timestamp.
// Audit:
//   - Human-readable and machine-parseable log format.
// -------------------------------------------------------
func logMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[INFO] %s HTTP %s %s\n", utcNow(), r.Method, r.URL.Path)
        h.ServeHTTP(w, r)
    })
}
