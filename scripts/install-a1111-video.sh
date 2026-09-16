#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEBUI="${WUJI_A1111_DIR:-$ROOT/../../stable-diffusion-webui}"
EXT_DIR="$WEBUI/extensions/sd-webui-text2video"
MODEL_DIR="$WEBUI/models/text2video/t2v"
HF_REPO="https://huggingface.co/kabachuha/modelscope-damo-text2video-pruned-weights/resolve/main"

if [[ ! -f "$WEBUI/webui.sh" ]]; then
  echo "A1111 not found at $WEBUI — run scripts/install-a1111.sh first."
  exit 1
fi

if [[ ! -d "$EXT_DIR/.git" ]]; then
  echo "Installing sd-webui-text2video extension ..."
  git clone --depth 1 https://github.com/kabachuha/sd-webui-text2video.git "$EXT_DIR"
else
  echo "Extension already present: $EXT_DIR"
fi

"$ROOT/scripts/patch-a1111-text2video.sh"

# Legacy path from early installs — keep a symlink for older docs/tools.
if [[ -d "$WEBUI/models/ModelScope/t2v" && ! -e "$MODEL_DIR" ]]; then
  mkdir -p "$(dirname "$MODEL_DIR")"
  ln -sfn ../ModelScope/t2v "$MODEL_DIR"
fi

mkdir -p "$MODEL_DIR"
for file in configuration.json open_clip_pytorch_model.bin text2video_pytorch_model.pth VQGAN_autoencoder.pth; do
  dest="$MODEL_DIR/$file"
  if [[ -f "$dest" ]]; then
    echo "Model file already present: $dest"
    continue
  fi
  echo "Downloading $file ..."
  curl -L --fail --progress-bar -o "$dest" "$HF_REPO/$file"
done

echo ""
echo "Done. Restart A1111 (or let wuji auto-start it), then:"
echo "  $ROOT/bin/wuji video \"a cat walking on a beach\" --duration 2 --fps 15"
echo ""
echo "Tip: in A1111 Settings → text2video, set 'Keep model in VRAM' to All for faster repeat runs."
