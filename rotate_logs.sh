#!/usr/bin/env bash
#-------------------------------------------------------
# usr/local/bin/rotate_logs.sh
#-------------------------------------------------------
# Purpose Summary:
#   - Rotate, hash, and archive audit logs under /evidence/.
#   - Enforce PNCRL retention rules:
#       * Logs: 180 days
#       * Hashes: 365 days
# Audit:
#   - Each action writes UTC ISO 8601 timestamped lines.
#   - All operations logged to /evidence/logs/rotation.log.
#   - Fails fast and reports any missing or invalid directories.
#-------------------------------------------------------

set -eu

LOG_DIR="/evidence/logs"
HASH_DIR="/evidence/hashes"
DATE_UTC="$(date -u +'%Y%m%d')"
NOW="$(date -u +'%FT%TZ')"
ROTATION_LOG="${LOG_DIR}/rotation.log"

mkdir -p "${LOG_DIR}" "${HASH_DIR}"

echo "[INFO] ${NOW} Starting daily rotation job" >> "${ROTATION_LOG}"

#-------------------------------------------------------
# Step 1: Locate today's log file
#-------------------------------------------------------
TARGET_LOG="${LOG_DIR}/requests_${DATE_UTC}.log"
if [ ! -f "${TARGET_LOG}" ] || [ ! -s "${TARGET_LOG}" ]; then
    echo "[INFO] ${NOW} No active log found for ${DATE_UTC}" >> "${ROTATION_LOG}"
    exit 0
fi

#-------------------------------------------------------
# Step 2: Compress and hash log file
#-------------------------------------------------------
GZ_FILE="${TARGET_LOG}.gz"
HASH_FILE="${HASH_DIR}/requests_${DATE_UTC}.sha512"

if gzip -c "${TARGET_LOG}" > "${GZ_FILE}"; then
    sha512sum "${GZ_FILE}" > "${HASH_FILE}"
    rm -f "${TARGET_LOG}"
    echo "[INFO] ${NOW} Archived ${TARGET_LOG} -> ${GZ_FILE}" >> "${ROTATION_LOG}"
    echo "[INFO] ${NOW} Created hash ${HASH_FILE}" >> "${ROTATION_LOG}"
else
    echo "[ERROR] ${NOW} Compression failed for ${TARGET_LOG}" >> "${ROTATION_LOG}"
    exit 1
fi

#-------------------------------------------------------
# Step 3: Apply retention rules
#-------------------------------------------------------
find "${LOG_DIR}" -type f -name 'requests_*.gz' -mtime +180 -print -delete >> "${ROTATION_LOG}" 2>&1
find "${HASH_DIR}" -type f -name 'requests_*.sha512' -mtime +365 -print -delete >> "${ROTATION_LOG}" 2>&1

echo "[INFO] ${NOW} Rotation completed successfully" >> "${ROTATION_LOG}"
exit 0
