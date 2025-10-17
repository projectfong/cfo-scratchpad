// -------------------------------------------------------
// backend/main.go
// -------------------------------------------------------
// Purpose Summary:
//   - Entry point for cfo-scratchpad backend service.
//   - Initializes secure REST API routes for folder and file handling.
//   - Serves static frontend assets from ./frontend via HTTP root path.
// Audit:
//   - Logs all actions with UTC ISO 8601 timestamps.
//   - Fails fast on any binding or dependency error.
//   - All actions are logged at appropriate levels: [INFO], [ERROR], [DEBUG].
//   - Exposes no unsafe routes or file traversal risks.
// -------------------------------------------------------

package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "cfo-scratchpad/handlers"
)

const (
    defaultPort   = "8080"
    staticDirPath = "./frontend"
)

// -------------------------------------------------------
// func utcNow()
// -------------------------------------------------------
// Purpose:
//   - Returns the current UTC timestamp in ISO 8601 format.
// Audit:
//   - Ensures all log timestamps are consistent and auditable.
// -------------------------------------------------------
func utcNow() string {
    return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

// -------------------------------------------------------
// func logInfo()
// -------------------------------------------------------
// Purpose:
//   - Logs informational messages with UTC timestamp.
// Audit:
//   - Ensures traceability of normal operations.
// -------------------------------------------------------
func logInfo(message string) {
    log.Printf("[INFO] %s %s\n", utcNow(), message)
}

// -------------------------------------------------------
// func logError()
// -------------------------------------------------------
// Purpose:
//   - Logs error messages with UTC timestamp.
// Audit:
//   - Ensures failures are clearly visible in logs.
// -------------------------------------------------------
func logError(message string) {
    log.Printf("[ERROR] %s %s\n", utcNow(), message)
}

// -------------------------------------------------------
// func main()
// -------------------------------------------------------
// Purpose:
//   - Initializes all HTTP routes and starts backend server.
//   - Handles API endpoints and serves static frontend files.
// Audit:
//   - Logs all startup actions and fails on port bind errors.
//   - Ensures all handlers are secure, minimal, and logged.
// -------------------------------------------------------
func main() {
    mux := http.NewServeMux()

    // API routes
    mux.HandleFunc("/folders", handlers.HandleFolders)
    mux.HandleFunc("/files", handlers.HandleFileList)
    mux.HandleFunc("/file", handlers.HandleFileGet)
    mux.HandleFunc("/file/save", handlers.HandleFileSave)
    mux.HandleFunc("/file/move", handlers.HandleFileMove)

    // Static frontend
    fs := http.FileServer(http.Dir(staticDirPath))
    mux.Handle("/", fs)

    port := os.Getenv("PORT")
    if port == "" {
        port = defaultPort
    }

    logInfo("Binding routes and starting server on port " + port)

    // Wrap all routes in AuditMiddleware to capture request evidence
    auditedMux := AuditMiddleware(mux) 

    err := http.ListenAndServe(":"+port, auditedMux)
    if err != nil {
        logError("Server failed to start: " + err.Error())
        os.Exit(1)
    }
}
