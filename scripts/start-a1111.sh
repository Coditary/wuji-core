#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEBUI="${WUJI_A1111_DIR:-$ROOT/../../stable-diffusion-webui}"
API_PORT="${WUJI_A1111_PORT:-7860}"
API_BASE="http://127.0.0.1:${API_PORT}"

if [[ ! -f "$WEBUI/webui.sh" ]]; then
  echo "A1111 not found at $WEBUI" >&2
  echo "Run scripts/install-a1111.sh first." >&2
  exit 1
fi

if curl -sf "${API_BASE}/sdapi/v1/options" >/dev/null 2>&1; then
  echo "A1111 API already running at ${API_BASE}" >&2
  echo "Stop it first if you need a fresh start: fuser -k ${API_PORT}/tcp" >&2
  exit 0
fi

cd "$WEBUI"
exec ./webui.sh "$@"
