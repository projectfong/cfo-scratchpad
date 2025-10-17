# -------------------------------------------------------
# Makefile
# -------------------------------------------------------
# Purpose Summary:
#   - Automate local build, run, stop, and clean actions.
# Audit:
#   - All actions log visibly. Fails fast on error.
# -------------------------------------------------------

APP = cfo-scratchpad

.PHONY: all build up down clean logs

all: build

# -------------------------------------------------------
# build
# -------------------------------------------------------
# Purpose:
#   - Builds Docker image from scratch.
# Audit:
#   - Uses explicit Dockerfile. Logs on success/failure.
# -------------------------------------------------------
build:
	@echo "[INFO] $(shell date -u +%FT%TZ) Building $(APP) Docker image..."
	docker compose build

# -------------------------------------------------------
# up
# -------------------------------------------------------
# Purpose:
#   - Starts scratchpad container using docker-compose.
# Audit:
#   - Container runs in background unless stopped.
# -------------------------------------------------------
up:
	@echo "[INFO] $(shell date -u +%FT%TZ) Starting $(APP)..."
	docker compose up -d

# -------------------------------------------------------
# down
# -------------------------------------------------------
# Purpose:
#   - Stops all running containers.
# Audit:
#   - Cleanly shuts down services.
# -------------------------------------------------------
down:
	@echo "[INFO] $(shell date -u +%FT%TZ) Stopping $(APP)..."
	docker compose down

# -------------------------------------------------------
# clean
# -------------------------------------------------------
# Purpose:
#   - Removes containers and local build artifacts.
# Audit:
#   - Removes all docker volumes and images.
# -------------------------------------------------------
clean:
	@echo "[WARN] $(shell date -u +%FT%TZ) Full clean of all data and images..."
	docker compose down --volumes --remove-orphans
	docker rmi $$(docker images -q $(APP)) || true

# -------------------------------------------------------
# logs
# -------------------------------------------------------
# Purpose:
#   - Follows backend logs in real-time.
# Audit:
#   - Includes timestamps and real-time events.
# -------------------------------------------------------
logs:
	docker compose logs -f
