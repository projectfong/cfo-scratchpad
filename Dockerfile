# -------------------------------------------------------
# Dockerfile
# -------------------------------------------------------
# Purpose:
#   - Build static Go backend and run on Alpine.
#   - Create entrypoint inside image with Unix LF endings.
#   - Fix /scratchpad-data ownership at startup, then drop to appuser.
# Audit:
#   - Deterministic build; startup emits UTC ISO 8601 logs.
#   - Avoids "exec format error" by generating script in-image.
# -------------------------------------------------------

# -------------------------------------------------------
# Build Stage
# -------------------------------------------------------
FROM golang:1.21 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /app

COPY . .

WORKDIR /app/backend
RUN printf "module cfo-scratchpad\n\ngo 1.21\n" > go.mod \
 && go mod tidy

RUN date -u +"[INFO] %Y-%m-%dT%H:%M:%SZ Building static binary" \
 && GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /cfo-scratchpad .

# -------------------------------------------------------
# Runtime Stage
# -------------------------------------------------------
FROM alpine:latest

# -------------------------------------------------------
# System Setup
# -------------------------------------------------------
# su-exec: safe privilege dropper
# dcron: lightweight cron daemon for daily rotation
RUN apk add --no-cache su-exec dcron

# Create runtime user
RUN adduser -D appuser

WORKDIR /home/appuser

# Backend binary and frontend
COPY --from=builder /cfo-scratchpad .
COPY --from=builder /app/frontend ./frontend

# -------------------------------------------------------
# Create entrypoint in-image (LF endings guaranteed)
# -------------------------------------------------------
RUN cat <<'SH' > /usr/local/bin/entrypoint.sh
#!/bin/sh
set -eu

# -------------------------------------------------------
# /usr/local/bin/entrypoint.sh
# -------------------------------------------------------
# Purpose Summary:
#   - Ensure /scratchpad-data exists and is owned by appuser:appuser.
#   - Start cron service for daily evidence rotation.
#   - Drop privileges to appuser and exec backend.
# Audit:
#   - All actions use UTC ISO 8601 timestamps.
#   - Fails fast on any error; no silent fallbacks.
# -------------------------------------------------------

utc_now() { date -u +"%Y-%m-%dT%H:%M:%SZ"; }
log_info()  { printf "[INFO]  %s %s\n" "$(utc_now)" "$*"; }
log_error() { printf "[ERROR] %s %s\n" "$(utc_now)" "$*" >&2; }

ROOT="/scratchpad-data"
ROTATION_LOG="/evidence/logs/rotation.log"

# Step 1: Prepare scratchpad directory
if [ ! -d "$ROOT" ]; then
    log_info "Creating scratch root: $ROOT"
    if ! mkdir -p "$ROOT"; then
        log_error "Failed to create $ROOT"
        exit 1
    fi
fi

# Step 2: Ensure ownership
if ! chown -R appuser:appuser "$ROOT" 2>/dev/null; then
    log_error "chown failed on $ROOT (check bind mount permissions on host)"
    exit 1
fi
log_info "Ownership set for $ROOT (appuser:appuser)"

# Step 3: Start cron for evidence rotation
mkdir -p "$(dirname "$ROTATION_LOG")"
if crond -b -L "$ROTATION_LOG" -c /etc/crontabs -p /tmp/crond.pid; then
    log_info "Cron service started for daily log rotation"
else
    log_error "Failed to start cron service; rotation may not occur"
fi

# Step 4: Launch backend as appuser
log_info "Starting backend as appuser"
exec su-exec appuser:appuser ./cfo-scratchpad
SH
RUN chmod 0755 /usr/local/bin/entrypoint.sh

# -------------------------------------------------------
# Install and schedule daily rotation (00:00 UTC)
# -------------------------------------------------------
COPY rotate_logs.sh /usr/local/bin/rotate_logs.sh
RUN chmod 755 /usr/local/bin/rotate_logs.sh \
 && echo "0 0 * * * /usr/local/bin/rotate_logs.sh" > /etc/crontabs/root

# -------------------------------------------------------
# Runtime wiring
# -------------------------------------------------------
RUN mkdir -p /evidence/logs && chown -R appuser:appuser /evidence
VOLUME /scratchpad-data
VOLUME /evidence
VOLUME /tmp
EXPOSE 8888
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
