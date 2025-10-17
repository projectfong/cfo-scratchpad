// -------------------------------------------------------
// backend/handlers/folders.go
// -------------------------------------------------------
// Purpose Summary:
//   - Handle folder listing and creation for cfo-scratchpad.
//   - Responds to GET (list) and POST (create) requests on /folders.
// Audit:
//   - Logs all operations with UTC ISO 8601 timestamps.
//   - Enforces path safety and fails fast on invalid input.
// -------------------------------------------------------

package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const scratchRoot = "/scratchpad-data"

// -------------------------------------------------------
// func utcNow()
// -------------------------------------------------------
// Purpose:
//   - Returns current UTC time in ISO 8601 format.
// Audit:
//   - Centralizes timestamp formatting for all logs.
// -------------------------------------------------------
func utcNow() string {
    return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

// -------------------------------------------------------
// func logInfo()
// -------------------------------------------------------
// Purpose:
//   - Log informational events.
// Audit:
//   - Visible, timestamped log trail for normal operations.
// -------------------------------------------------------
func logInfo(msg string) {
    fmt.Printf("[INFO] %s %s\n", utcNow(), msg)
}

// -------------------------------------------------------
// func logError()
// -------------------------------------------------------
// Purpose:
//   - Log operational errors.
// Audit:
//   - Ensures all failure paths are recorded visibly.
// -------------------------------------------------------
func logError(msg string) {
    fmt.Printf("[ERROR] %s %s\n", utcNow(), msg)
}

// -------------------------------------------------------
// func sanitizePath()
// -------------------------------------------------------
// Purpose:
//   - Prevents directory traversal by sanitizing input paths.
// Audit:
//   - Strips `..` and ensures paths are rooted under scratchRoot.
// -------------------------------------------------------
func sanitizePath(path string) string {
    clean := filepath.Clean(path)
    if strings.Contains(clean, "..") {
        return ""
    }
    return filepath.Join(scratchRoot, clean)
}

// -------------------------------------------------------
// func HandleFolders(w http.ResponseWriter, r *http.Request)
// -------------------------------------------------------
// Purpose:
//   - Dispatch handler for GET (list folders) and POST (create folder).
// Audit:
//   - Logs method, path, and outcomes for all folder actions.
// -------------------------------------------------------
func HandleFolders(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        handleListFolders(w, r)
    case "POST":
        handleCreateFolder(w, r)
    default:
        logError("Unsupported method: " + r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// -------------------------------------------------------
// func handleListFolders(w, r)
// -------------------------------------------------------
// Purpose:
//   - Lists all subfolders under the scratchpad root.
// Audit:
//   - Logs total folders found and any filesystem errors.
//   - Ensures JSON response is always an array (never null).
//   - UTC ISO 8601 timestamps via logInfo/logError.
// -------------------------------------------------------
func handleListFolders(w http.ResponseWriter, r *http.Request) {
    // Always initialize to an empty slice so JSON is [] instead of null.
    folders := []string{}

    // If root is missing, treat as empty but log clearly.
    if _, statErr := os.Stat(scratchRoot); os.IsNotExist(statErr) {
        logInfo("Scratch root missing; returning empty folder list: " + scratchRoot)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(folders)
        return
    }

    err := filepath.Walk(scratchRoot, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && path != scratchRoot {
            rel, relErr := filepath.Rel(scratchRoot, path)
            if relErr != nil {
                return relErr
            }
            folders = append(folders, rel)
        }
        return nil
    })

    if err != nil {
        logError("Failed to list folders: " + err.Error())
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    logInfo(fmt.Sprintf("Listed %d folders", len(folders)))

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(folders)
}

// -------------------------------------------------------
// func handleCreateFolder(w, r)
// -------------------------------------------------------
// Purpose:
//   - Creates a new folder under scratchpad root.
// Audit:
//   - Logs created path and fails fast on unsafe paths.
// -------------------------------------------------------
func handleCreateFolder(w http.ResponseWriter, r *http.Request) {
    type Request struct {
        Name string `json:"name"`
    }

    var req Request
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Name == "" {
        logError("Invalid folder creation payload")
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    safePath := sanitizePath(req.Name)
    if safePath == "" {
        logError("Rejected unsafe folder name: " + req.Name)
        http.Error(w, "Invalid folder path", http.StatusBadRequest)
        return
    }

    mkErr := os.MkdirAll(safePath, 0755)
    if mkErr != nil {
        logError("Failed to create folder: " + mkErr.Error())
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    logInfo("Created folder: " + safePath)
    w.WriteHeader(http.StatusCreated)
}
