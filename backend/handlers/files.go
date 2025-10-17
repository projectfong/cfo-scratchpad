// -------------------------------------------------------
// backend/handlers/files.go
// -------------------------------------------------------
// Purpose Summary:
//   - Handle file list, read, write, and move for .txt files.
// Audit:
//   - Returns JSON arrays (never null). Logs with UTC ISO 8601.
//   - Fails fast with clear HTTP status codes.
// -------------------------------------------------------

package handlers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

const fileExt = ".txt"

// -------------------------------------------------------
// func HandleFileList(w, r)
// -------------------------------------------------------
// Purpose:
//   - List .txt files in a sanitized folder under scratchpad root.
// Audit:
//   - Always JSON encodes an array ([] when empty).
//   - Logs counts and errors with UTC ISO 8601 timestamps.
// -------------------------------------------------------
func HandleFileList(w http.ResponseWriter, r *http.Request) {
    folder := r.URL.Query().Get("folder")
    absPath := sanitizePath(folder)

    // Always start with an initialized slice so JSON is [] not null.
    files := []string{}

    if absPath == "" {
        logError("Invalid folder path requested: " + folder)
        http.Error(w, "Invalid folder path", http.StatusBadRequest)
        return
    }

    // If the folder does not exist, treat as empty list.
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        logInfo("Folder does not exist; returning empty list: " + absPath)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(files)
        return
    }

    entries, err := ioutil.ReadDir(absPath)
    if err != nil {
        logError("Failed to read folder: " + err.Error())
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), fileExt) {
            files = append(files, entry.Name())
        }
    }

    logInfo(fmt.Sprintf("Listed %d files in folder: %s", len(files), absPath))
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}

// -------------------------------------------------------
// func HandleFileGet(w, r)
// -------------------------------------------------------
// Purpose:
//   - Returns the contents of a specific .txt file under scratchpad root.
// Audit:
//   - Logs path read and any read failures with UTC ISO 8601 timestamps.
// -------------------------------------------------------
func HandleFileGet(w http.ResponseWriter, r *http.Request) {
    file := r.URL.Query().Get("path")
    absPath := sanitizePath(file)

    if absPath == "" || !strings.HasSuffix(absPath, fileExt) {
        logError("Invalid file path requested: " + file)
        http.Error(w, "Invalid file path", http.StatusBadRequest)
        return
    }

    content, err := ioutil.ReadFile(absPath)
    if err != nil {
        logError("Failed to read file: " + absPath + " - " + err.Error())
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    logInfo("Read file: " + absPath)

    w.Header().Set("Content-Type", "text/plain")
    w.Write(content)
}

// -------------------------------------------------------
// func HandleFileSave(w, r)
// -------------------------------------------------------
// Purpose:
//   - Saves or updates content to a specific .txt file.
// Audit:
//   - Logs before/after snapshot of saved file (truncated for safety).
//   - Sanitizes paths and logs full path written to with UTC timestamps.
// -------------------------------------------------------
func HandleFileSave(w http.ResponseWriter, r *http.Request) {
    type SaveRequest struct {
        Path    string `json:"path"`
        Content string `json:"content"`
    }

    var req SaveRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Path == "" {
        logError("Invalid save request payload")
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    absPath := sanitizePath(req.Path)
    if absPath == "" || !strings.HasSuffix(absPath, fileExt) {
        logError("Rejected unsafe save path: " + req.Path)
        http.Error(w, "Invalid file path", http.StatusBadRequest)
        return
    }

    before := ""
    if existing, readErr := ioutil.ReadFile(absPath); readErr == nil {
        before = string(existing)
    }

    err = ioutil.WriteFile(absPath, []byte(req.Content), 0644)
    if err != nil {
        logError("Failed to save file: " + absPath + " - " + err.Error())
        http.Error(w, "Write failed", http.StatusInternalServerError)
        return
    }

    logInfo("Saved file: " + absPath)
    logInfo("Before snapshot: " + truncateLog(before))
    logInfo("After snapshot: " + truncateLog(req.Content))

    w.WriteHeader(http.StatusOK)
}

// -------------------------------------------------------
// func HandleFileMove(w, r)
// -------------------------------------------------------
// Purpose:
//   - Moves a file from one folder to another safely.
// Audit:
//   - Logs full old/new paths and fails fast on any invalid input.
//   - UTC ISO 8601 timestamps via logInfo/logError.
// -------------------------------------------------------
func HandleFileMove(w http.ResponseWriter, r *http.Request) {
    type MoveRequest struct {
        From string `json:"from"`
        To   string `json:"to"`
    }

    var req MoveRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.From == "" || req.To == "" {
        logError("Invalid move request payload")
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    fromPath := sanitizePath(req.From)
    toPath := sanitizePath(req.To)

    if fromPath == "" || toPath == "" || !strings.HasSuffix(fromPath, fileExt) || !strings.HasSuffix(toPath, fileExt) {
        logError("Rejected unsafe move paths: " + req.From + " -> " + req.To)
        http.Error(w, "Invalid file paths", http.StatusBadRequest)
        return
    }

    err = os.Rename(fromPath, toPath)
    if err != nil {
        logError("Failed to move file: " + err.Error())
        http.Error(w, "Move failed", http.StatusInternalServerError)
        return
    }

    logInfo("Moved file: " + fromPath + " -> " + toPath)
    w.WriteHeader(http.StatusOK)
}

// -------------------------------------------------------
// func truncateLog(text string) string
// -------------------------------------------------------
// Purpose:
//   - Truncates long content logs for readability and safety.
// Audit:
//   - Prevents logging of large/malformed files that clutter logs.
//   - Logs remain human-readable and machine-parseable.
// -------------------------------------------------------
func truncateLog(text string) string {
    max := 200
    if len(text) <= max {
        return text
    }
    return text[:max] + " ... [truncated]"
}
