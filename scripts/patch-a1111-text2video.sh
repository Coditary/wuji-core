#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEBUI="${WUJI_A1111_DIR:-$ROOT/../../stable-diffusion-webui}"
EXT="$WEBUI/extensions/sd-webui-text2video"

if [[ ! -d "$EXT/scripts" ]]; then
  echo "text2video extension not found at $EXT — run scripts/install-a1111-video.sh first."
  exit 1
fi

python3 - "$EXT" <<'PY'
import sys
from pathlib import Path

ext = Path(sys.argv[1])

p1 = ext / "scripts/modelscope/process_modelscope.py"
text = p1.read_text()
if "if len(stable_lora_args) >= 8:" not in text:
    text = text.replace(
        "    stable_lora_args = stable_lora_processor.process_extension_args(all_args=extra_args) \n    stable_lora_processor.process(pipe, *stable_lora_args)",
        "    stable_lora_args = stable_lora_processor.process_extension_args(all_args=extra_args)\n    if len(stable_lora_args) >= 8:\n        stable_lora_processor.process(pipe, *stable_lora_args)",
    )
    text = text.replace(
        "    stable_lora_args = stable_lora_processor.process_extension_args(all_args=extra_args)\n    stable_lora_processor.process(pipe, *stable_lora_args)",
        "    stable_lora_args = stable_lora_processor.process_extension_args(all_args=extra_args)\n    if len(stable_lora_args) >= 8:\n        stable_lora_processor.process(pipe, *stable_lora_args)",
    )
    p1.write_text(text)
    print("patched process_modelscope.py")

p2 = ext / "scripts/modelscope/clip_hardcode.py"
text2 = p2.read_text()
old = "        if opts.enable_emphasis:"
new = "        if getattr(opts, 'enable_emphasis', True):"
if old in text2:
    p2.write_text(text2.replace(old, new, 1))
    print("patched clip_hardcode.py")

print("text2video compatibility patches applied")
PY
