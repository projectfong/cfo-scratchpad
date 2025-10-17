//-------------------------------------------------------
// backend/middleware_audit.go
//-------------------------------------------------------
// Purpose Summary:
//   - Record all file and API access events for audit evidence.
//   - Maintain immutable event trail under /evidence/logs/.
// Audit:
//   - Appends JSON records to /evidence/logs/requests_YYYY-MM-DD.log.
//   - Emits UTC ISO 8601 timestamps for every action and error.
//   - Never creates or modifies directories.
//   - Fails safe if /evidence/logs/ is missing or unwritable.
// Compliance:
//   - Required under PNCRL-AUDIT-1.0 non-commercial license terms.
//   - Evidence logs must be retained and hashed per rotation policy.
//-------------------------------------------------------

package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"
)

//-------------------------------------------------------
// Struct: AuditEvent
//-------------------------------------------------------
// Purpose:
//   - Define a consistent JSON schema for auditable events.
// Audit:
//   - One event per line in /evidence/logs/requests_YYYY-MM-DD.log.
//   - Immutable once written (append-only).
//-------------------------------------------------------
type AuditEvent struct {
    Timestamp string `json:"timestamp"`
    Method    string `json:"method"`
    Path      string `json:"path"`
    RemoteIP  string `json:"remote_ip"`
    Status    int    `json:"status"`
    Duration  int64  `json:"duration_ms"`
}

//-------------------------------------------------------
// Function: AuditMiddleware
//-------------------------------------------------------
// Purpose:
//   - Wrap HTTP handlers to capture metadata on every request.
// Audit:
//   - Captures method, path, remote IP, response code, and latency.
//   - Delegates event persistence to writeAuditEvent().
//   - Emits one structured JSON audit record per request.
//-------------------------------------------------------
func AuditMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now().UTC()
        lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(lrw, r)

        event := AuditEvent{
            Timestamp: start.Format(time.RFC3339), // ISO 8601 UTC timestamp
            Method:    r.Method,
            Path:      r.URL.Path,
            RemoteIP:  r.RemoteAddr,
            Status:    lrw.statusCode,
            Duration:  time.Since(start).Milliseconds(),
        }

        writeAuditEvent(event)
    })
}

//-------------------------------------------------------
// Struct: loggingResponseWriter
//-------------------------------------------------------
// Purpose:
//   - Capture the final HTTP status code from handler responses.
// Audit:
//   - Ensures status codes are correctly logged in each event.
//-------------------------------------------------------
type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

//-------------------------------------------------------
// Function: writeAuditEvent
//-------------------------------------------------------
// Purpose:
//   - Append JSON audit events to daily evidence logs.
// Audit:
//   - File naming: /evidence/logs/requests_YYYY-MM-DD.log
//   - Creates the log file if missing (license-permitted).
//   - Never creates directories; /evidence/logs must pre-exist.
//   - Each JSON record represents one auditable transaction.
//   - Logs [ERROR] with UTC ISO 8601 timestamp on any failure.
//-------------------------------------------------------
func writeAuditEvent(event AuditEvent) {
    logDir := "/evidence/logs"
    logFile := filepath.Join(logDir, "requests_"+time.Now().UTC().Format("2006-01-02")+".log")

    // Verify that /evidence/logs directory exists and is valid
    if stat, err := os.Stat(logDir); err != nil || !stat.IsDir() {
        log.Printf("[ERROR] %s audit path missing or invalid: %s (%v)",
            time.Now().UTC().Format(time.RFC3339), logDir, err)
        return
    }

    // Open or create the daily log file for appending
    f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("[ERROR] %s audit open failed: %v",
            time.Now().UTC().Format(time.RFC3339), err)
        return
    }
    defer f.Close()

    // Encode event as a single JSON line
    enc := json.NewEncoder(f)
    if err := enc.Encode(event); err != nil {
        log.Printf("[ERROR] %s audit encode failed: %v",
            time.Now().UTC().Format(time.RFC3339), err)
    }
}
