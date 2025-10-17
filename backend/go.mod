// -------------------------------------------------------
// backend/go.mod
// -------------------------------------------------------
// Purpose Summary:
//   - Defines the Go module name for cfo-scratchpad backend.
//   - Enables correct internal import paths and module resolution.
// Audit:
//   - Must match all import paths in source files (e.g., cfo-scratchpad/handlers).
//   - Required for reproducible, containerized builds.
//   - Regenerated with `go mod tidy` when dependencies change.
// -------------------------------------------------------

module cfo-scratchpad

go 1.21
