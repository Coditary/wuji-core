#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEBUI="${WUJI_A1111_DIR:-$ROOT/../../stable-diffusion-webui}"
MODEL_DIR="$WEBUI/models/Stable-diffusion"
MODEL_FILE="sd-v1-5-inpainting.ckpt"
MODEL_URL="https://huggingface.co/runwayml/stable-diffusion-inpainting/resolve/main/sd-v1-5-inpainting.ckpt"

if [[ ! -f "$WEBUI/webui.sh" ]]; then
  echo "A1111 not found at $WEBUI — run scripts/install-a1111.sh first."
  exit 1
fi

mkdir -p "$MODEL_DIR"
dest="$MODEL_DIR/$MODEL_FILE"
if [[ -f "$dest" ]]; then
  echo "Inpaint model already present: $dest"
else
  echo "Downloading SD 1.5 inpainting checkpoint (~4 GB) ..."
  curl -L --fail --progress-bar -o "$dest" "$MODEL_URL"
fi

echo ""
echo "Done. Set as default inpaint model:"
echo "  $ROOT/bin/wuji config set a1111 default_inpaint_model $MODEL_FILE"
echo "Then retry inpainting (driver will auto-select this model for --task inpaint)."
