# -------------------------------------------------------
# entrypoint.sh
# -------------------------------------------------------
# Purpose Summary:
#   - Ensure /scratchpad-data and /evidence/logs exist and are owned by appuser:appuser.
#   - Initialize evidence rotation scheduler (cron) inside container.
#   - Drop privileges to appuser and execute backend binary.
# Audit:
#   - Prints all actions with UTC ISO 8601 timestamps and levels.
#   - Logs cron startup and backend launch.
#   - Fails fast on any error; no silent fallbacks.
# -------------------------------------------------------

#!/bin/sh
set -eu

# -------------------------------------------------------
# Function: utc_now
# -------------------------------------------------------
utc_now() { date -u +"%Y-%m-%dT%H:%M:%SZ"; }

# -------------------------------------------------------
# Logging helpers
# -------------------------------------------------------
log_info()  { printf "[INFO]  %s %s\n" "$(utc_now)" "$*"; }
log_error() { printf "[ERROR] %s %s\n" "$(utc_now)" "$*" >&2; }

ROOT="/scratchpad-data"
EVIDENCE_DIR="/evidence/logs"
ROTATION_LOG="${EVIDENCE_DIR}/rotation.log"

# -------------------------------------------------------
# Step 1: Verify and prepare scratchpad directory
# -------------------------------------------------------
if [ ! -d "$ROOT" ]; then
    log_info "Creating scratch root: $ROOT"
    if ! mkdir -p "$ROOT"; then
        log_error "Failed to create $ROOT"
        exit 1
    fi
fi

# -------------------------------------------------------
# Step 2: Ensure /evidence/logs exists and ownership is correct
# -------------------------------------------------------
if [ ! -d "$EVIDENCE_DIR" ]; then
    log_info "Creating evidence log directory: $EVIDENCE_DIR"
    if ! mkdir -p "$EVIDENCE_DIR"; then
        log_error "Failed to create $EVIDENCE_DIR"
        exit 1
    fi
fi

# Force ownership correction on mounted volumes
if chown -R appuser:appuser "$ROOT" "$EVIDENCE_DIR" 2>/dev/null; then
    log_info "Ownership set for $ROOT and $EVIDENCE_DIR (appuser:appuser)"
else
    log_error "chown failed on one or more mount paths (check Docker volume or host permissions)"
    ls -ld "$ROOT" "$EVIDENCE_DIR" || true
    exit 1
fi

# -------------------------------------------------------
# Step 3: Start cron for daily evidence rotation
# -------------------------------------------------------
log_info "Starting cron service for daily log rotation"
if crond -b -L "$ROTATION_LOG" -c /etc/crontabs -p /tmp/crond.pid; then
    log_info "Cron service started successfully"
else
    log_error "Failed to start cron service; rotation may not occur"
fi

# -------------------------------------------------------
# Step 4: Launch backend as non-root user
# -------------------------------------------------------
log_info "Starting backend as appuser"
exec su-exec appuser:appuser ./cfo-scratchpad
