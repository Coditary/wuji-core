#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEBUI="${WUJI_A1111_DIR:-$ROOT/../../stable-diffusion-webui}"
MODEL_DIR="$WEBUI/models/Stable-diffusion"
MODEL_FILE="v1-5-pruned-emaonly.safetensors"
MODEL_URL="https://huggingface.co/stable-diffusion-v1-5/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.safetensors"

if [[ ! -f "$WEBUI/webui.sh" ]]; then
  echo "Cloning Automatic1111 into $WEBUI ..."
  git clone --depth 1 https://github.com/AUTOMATIC1111/stable-diffusion-webui.git "$WEBUI"
fi

mkdir -p "$MODEL_DIR"
if [[ ! -f "$MODEL_DIR/$MODEL_FILE" ]]; then
  echo "Downloading SD 1.5 checkpoint (~4 GB) ..."
  curl -L --fail --progress-bar -o "$MODEL_DIR/$MODEL_FILE" "$MODEL_URL"
else
  echo "Model already present: $MODEL_DIR/$MODEL_FILE"
fi

echo "Installing Python dependencies with micromamba Python 3.10 (may take several minutes) ..."
cd "$WEBUI"
export WUJI_A1111_PYTHON="${WUJI_A1111_PYTHON:-$HOME/micromamba/envs/sd-webui/bin/python}"
./webui.sh --exit --skip-torch-cuda-test --skip-python-version-check

echo ""
echo "Done. Start with:"
echo "  $ROOT/scripts/start-a1111.sh"
echo "For video generation, also run:"
echo "  $ROOT/scripts/install-a1111-video.sh"
echo "Then in another terminal:"
echo "  $ROOT/bin/wuji-driver-a1111"
echo "  $ROOT/bin/wuji config set driver video a1111"
echo "  $ROOT/bin/wuji video \"your prompt\""
