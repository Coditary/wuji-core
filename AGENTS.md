# Wuji Core — Agent Guide

Wuji Core is the backend daemon and orchestration layer. The CLI frontend lives in `../wuji-ai/`.

## Quick orientation

| Path | Purpose |
|------|---------|
| `cmd/wuji-core/main.go` | Core daemon entry point |
| `cmd/wuji-driver-dummy/main.go` | Standalone gRPC dummy driver |
| `pkg/core/` | Orchestrator: registry, scheduler, remote connect |
| `pkg/driver/` | Driver interfaces, types, proto mapping, built-in drivers |
| `internal/server/grpc/` | Core gRPC server — driver registry + WujiCore API |
| `api/proto/v1/driver.proto` | gRPC API contract |
| `../driver/` | External driver plugins (sibling of wuji-core) |

## Architecture

```
wuji CLI (../wuji-ai) ──gRPC──► wuji-core daemon
                                    │
                                    ├── Scheduler (RAM/VRAM)
                                    └── Registry ──► Driver (built-in or gRPC)
                                                         ▲
                                    wuji-driver-* ───────┘
```

`wuji serve` runs the long-lived core on a unix socket (default: `.wuji/run/drivers/core.sock`).
Every other `wuji` command connects to that daemon (auto-starts it if missing).
Remote clients set `WUJI_CORE_ADDR=host:port` (TCP). Media inputs (files and directories, including
train datasets and LoRA weights) upload via `PutFiles` with zstd compression; outputs download via
`GetFiles` into `.wuji/remote-out/` (or the original export path for `wuji rag --export`). Set
`WUJI_SKIP_MEDIA_TRANSFER=1` to disable.
Drivers are loaded on demand — nothing starts at daemon boot except built-ins.
After `idle_shutdown_seconds` with no active wuji command (default 1800), remote drivers are stopped.
Active generate/train/load operations hold a busy lock for their full duration (hours-safe).
Set `drivers.resources.idle_shutdown_seconds: 0` to disable.

1. **Daemon**: `wuji serve` (or auto-start on first command) — central core on unix socket.
2. **Built-in**: `dummy`, `local`, `ffmpeg` registered in the daemon (lightweight).
3. **Remote drivers**: started lazily via `EnsureDriver` when a command needs them (unix socket per driver ID).
4. **Persisted**: `wuji driver connect <addr> --save` writes to `.wuji/config.yaml` (`driver_endpoints`).
5. **Optional push registration**: driver binaries with `--core <daemon-addr>`.
6. **Driver settings**: optional `.wuji/config.yaml`. Defaults work without any file.

## Capabilities

Defined in `internal/capability/capability.go`:

- `text`, `image`, `video`, `audio`, `mesh`, `voice`, `train`, `dataset`

Each driver advertises a subset. The core uses `driver.As[T](d, cap)` to call the right interface (`TextGenerator`, `ImageGenerator`, …).

## Drivers

| ID | Type | Location | Capabilities | Notes |
|----|------|----------|--------------|-------|
| `dummy` | built-in + gRPC | `internal/driver/dummy/`, `cmd/wuji-driver-dummy/` | all | Placeholder responses for every CLI command |
| `echo` | gRPC | `../driver/echo/` | `text` | Echoes prompts; no ML deps; default unix socket |
| `llama` | gRPC | `../driver/llama/` | `text` | llama.cpp or Ollama; config via `wuji config` |
| `vllm` | gRPC | `../driver/vllm/` | `text` | Starts `vllm serve` automatically; default unix socket |

### Building

```bash
make build          # wuji + wuji-driver-dummy → bin/
make build-drivers  # wuji-driver-llama + wuji-driver-echo + wuji-driver-vllm → bin/
make test
```

### Testing without ML backends

**Option A — built-in dummy (fastest)**

```bash
make build
./bin/wuji driver list
./bin/wuji generate text "hello"
./bin/wuji generate image "a cat" --width 256 --height 256
./bin/wuji text --audio meeting.wav
./bin/wuji dataset list
```

**Option B — remote dummy driver**

```bash
# terminal 1
./bin/wuji serve --addr 127.0.0.1:50051

# terminal 2
./bin/wuji-driver-dummy --core 127.0.0.1:50051

# terminal 3
./bin/wuji driver list          # shows dummy (built-in) + dummy (remote)
./bin/wuji generate text "hi" -d dummy
```

**Option C — echo driver (text only)**

```bash
make build-drivers
./bin/wuji-driver-echo
./bin/wuji driver connect unix://$PWD/.wuji/run/drivers/echo.sock
./bin/wuji generate text "ping" -d echo
```

**Option D — vLLM driver (starts vLLM server automatically)**

```bash
# optional: configure via Wuji (writes .wuji/config.yaml)
wuji config set vllm default_model meta-llama/Llama-3.2-1B
wuji config set vllm startup_timeout_seconds 600
wuji config list

make build-drivers
./bin/wuji-driver-vllm
./bin/wuji driver connect unix://$PWD/.wuji/run/drivers/vllm.sock
./bin/wuji generate text "Hallo" -d vllm
```

No per-driver config files under `../driver/vllm/` — defaults apply out of the box. Override via `wuji config set vllm <key> <value>`, CLI flags on `wuji-driver-vllm`, or env vars (`VLLM_MODEL`, `VLLM_BIN`, …).

## Writing a new driver

1. Implement `driver.Driver` (`Info`, `Capabilities`, `Close`).
2. Implement capability interfaces you support (e.g. `TextGenerator` with `GenerateText`).
3. Expose via gRPC using the shared host:

```go
grpcdriver.Serve(myDriver, grpcdriver.HostOptions{
    Addr:     "unix:///path/to/.wuji/run/drivers/my-driver.sock",
    CoreAddr: "", // or core address for auto-register
})
```

4. Add a `Makefile` under `../driver/<name>/` and register it in root `Makefile` `build-drivers`.

Each driver declares **capabilities** (text, image, …) and **supported source formats** per capability (gguf, safetensors, huggingface, ollama, …) in `Info().FormatSupport`. Remote drivers expose this via gRPC `DriverMetadata.format_support`.

```bash
wuji driver info llama   # after connecting the driver
wuji driver info vllm
wuji driver formats      # list of known format identifiers
```

## Config

- Project root: directory containing `go.mod` (or `WUJI_ROOT` env).
- `.wuji/config.yaml` — optional central settings (`default_driver`, `default_provider`, `providers`, `capability_drivers`, `driver_endpoints`, `drivers.vllm.*`).
- `.wuji/run/drivers/*.sock` — default local gRPC endpoints (unix sockets, one per driver ID).
- `.wuji/mcp.json` — MCP server definitions (Cursor-compatible `mcpServers` format).

```json
// .wuji/mcp.json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]
    },
    "exa": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "https://mcp.exa.ai/mcp"],
      "env": { "EXA_API_KEY": "your-key" }
    }
  }
}
```

### Declarative API drivers (`apis`)

Define remote HTTP backends in YAML — each entry becomes a **virtual driver** (no new Go code).

```yaml
apis:
  anthropic:
    text:
      protocol: anthropic
      api_key_env: ANTHROPIC_API_KEY
      model: claude-sonnet-4-20250514

  openai:
    text:
      protocol: openai
      api_key_env: OPENAI_API_KEY
      model: gpt-4o

  my-images:
    image:
      protocol: http
      method: POST
      url: https://api.example.com/v1/images
      headers:
        Authorization: "Bearer ${API_KEY}"
      body:
        prompt: "${prompt}"
        width: "${width}"
        height: "${height}"
      response:
        image_url: "$.data.url"
      api_key_env: EXAMPLE_API_KEY
```

Text protocols: `openai`, `anthropic`, `http`. Image: `openai-images`, `http`. Audio/video: `http`.  
Use: `wuji text "hi" --driver anthropic`. Agent clients use `GenerateText` with `messages` + `tools`.

#### Full `APICapabilitySpec` field reference

| Field | Description |
|-------|-------------|
| `protocol` | `openai`, `anthropic`, `openai-images`, `http` |
| `method`, `url`, `base_url` | HTTP method and endpoint |
| `query` | URL query params (template values) |
| `headers` | Extra request headers |
| `body` | JSON body key/value templates |
| `body_raw` | Free-form JSON body template (overrides `body`) |
| `content_type` | Override `Content-Type` (default `application/json`) |
| `response` | JSON paths: `text`, `image_url`, `image_base64`, `audio_url`, `video_url`, … |
| `error_path` | JSON path for error message extraction |
| `auth` | `{ type: bearer\|header\|query\|none, header, query_param, prefix }` |
| `defaults` | YAML defaults: `model`, `max_tokens`, `temperature`, `top_p`, `top_k`, `cfg_scale`, `steps`, `width`, `height`, … |
| `multipart` | `[{ name, value, file }]` for `multipart/form-data` |
| `poll` | Async jobs: `interval_seconds`, `timeout_seconds`, `job_id_path`, `poll_url`, `status_path`, `done_values`, `fail_values` |
| `streaming` | `true` for SSE text streaming (OpenAI/Anthropic delta formats) |
| `timeout_seconds` | HTTP client timeout |
| `api_key`, `api_key_env`, `model` | Shorthand auth/model (legacy-compatible) |
| `tasks` | Per-task overrides merged onto this spec (keys = driver task names) |

Template variables in `${name}`: request fields (`prompt`, `messages`, `width`, `height`, `task`, …) plus `API_KEY`, `model` from spec/defaults.

Capability blocks: `text`, `image`, `video`, `audio`, `mesh`, `voice`, `video2audio`, `train`, `dataset`, `data`, `rag`.  
Without `tasks:`, all tasks for that capability are advertised (use `${task}` in URL/body). With `tasks:`, only listed tasks are enabled.

Text multimodal tasks (under `text.tasks`): `image-to-text`, `video-to-text`, `audio-to-text`, `document-to-text`.

#### Template variables by capability

All blocks expose `${task}`, `${model}`, `${API_KEY}` (when configured). File paths in `multipart[].file` use the same names.

**text:** `prompt`, `system_prompt`, `input_mode`, `media_path`, `messages_json`, `tools_json`, `max_tokens`, `temperature`, `top_p`, `top_k`, `min_p`, `frequency_penalty`, `presence_penalty`, `repetition_penalty`, `stop_sequences_json`, `translate`, `target_lang`, `language`, `beam_size`, `word_timestamps`, `vad_enabled`, `vad_threshold`, `loras_json`, `seed`, `context_window`

**image:** `prompt`, `negative_prompt`, `width`, `height`, `steps`, `cfg_scale`, `sampler`, `batch_size`, `batch_count`, `seed`, `denoising_strength`, `init_image_path`, `mask_image_path`, `control_image_path`, `control_type`, `control_units_json`, `control_mode`, `style_image_path`, `style_weight`, `scale`, `loras_json`, `mode`, sprite fields (`frame_width`, `frame_height`, `columns`, `rows`, `sprite_action`, …)

**video:** `prompt`, `negative_prompt`, `duration`, `fps`, `frames`, `motion_strength`, `context_length`, `sampler`, `scheduler`, `camera_control`, `init_image_path`, `init_video_path`, `input_video_path`, `target_fps`, `scale`, `seed`

**audio:** `prompt`, `lyrics`, `negative_prompt`, `voice`, `language`, `duration`, `overlap`, `temperature`, `cfg_scale`, `top_p`, `top_k`, `sample_rate`, `reference_path`, `format`, `speed`, `pitch`, `emotion`, `style`, `energy`, `seed`

**mesh:** `prompt`, `format`, `target_tris`, `scale`, `mode`, `representation`, `init_mesh_path`, `high_mesh_path`, `init_image_path`, `init_image_paths_json`, `init_video_path`, `init_depth_path`, `init_point_cloud_path`, `init_splat_path`, `mask_path`, `style_image_path`, `texture_image_path`, `animation_path`

**voice:** `name`, `sample_path`, `source_path`, `target_model`, `pitch_shift`, `index_rate`, `protect`, `vocoder`, `denoise`, `denoise_strength`, `chunk_size`, `crossfade`

**data:** `text`, `texts_json`, `text_file`, `image_path`, `csv_path`, `graph_path`, `input_format`, `input_shape`, `output_shape`, `horizon`, `use_stdin`

**rag:** `store_root`, `collection`, `query`, `source_paths_json`, `chunk_size`, `chunk_overlap`, `embed_model`, `text_model`, `top_k`, `min_score`, `max_tokens`, `filter_json`, `index_mode`, `force`, `dry_run`, `glob`, `exclude_json`, `url`, `temperature`, `top_p`, `rename_to`, `export_path`, `import_path`, …

**dataset:** `name`, `path`, `dataset_id`, `source_path`, `recursive`, `format`, `description`, `tag`, `message`

**train:** per-capability vars (`dataset_id`, `base_model`, `epochs`, `learning_rate`, `lora_rank`, `batch_size`, `sample_path`, …) — see `internal/driver/api/vars_train.go`

**video2audio:** `video_path`

```yaml
apis:
  async-video:
    video:
      protocol: http
      method: POST
      url: https://api.example.com/v1/video
      api_key_env: EXAMPLE_API_KEY
      auth: { type: query, query_param: api_key }
      body: { prompt: "${prompt}", duration: "${duration}" }
      poll:
        job_id_path: "$.id"
        poll_url: "https://api.example.com/v1/jobs/${job_id}"
        status_path: "$.status"
        done_values: [completed, succeeded]
        fail_values: [failed, error]
        interval_seconds: 5
        timeout_seconds: 600
      response:
        video_url: "$.output.url"

  stream-chat:
    text:
      protocol: openai
      api_key_env: OPENAI_API_KEY
      streaming: true
      defaults: { max_tokens: 4096, temperature: 0.7 }
      timeout_seconds: 120
```

Legacy `providers:` entries are auto-migrated into `apis` (cloud types only). Taiji syncs driver ids via `ListProviders`.

### MCP server management

```bash
wuji mcp --import                          # import from .cursor/mcp.json or ~/.cursor/mcp.json
wuji mcp --import ~/.cursor/mcp.json       # import from specific file
wuji mcp add filesystem --command npx --arg -y --arg @modelcontextprotocol/server-filesystem --arg /tmp
wuji mcp list
wuji mcp start filesystem --stdio          # foreground MCP proxy (for OpenCode type: local)
wuji mcp start filesystem --http           # Streamable HTTP gateway (background)
wuji mcp start filesystem --sse --port 9333
wuji mcp start --all --http                # gateway per server (ports from --port-base)
wuji mcp sync opencode --write             # write opencode.json with wuji mcp start --stdio commands
wuji mcp sync opencode --http --write      # write remote URLs (start gateways first)
wuji mcp stop filesystem
wuji mcp stop --all
wuji mcp restart filesystem
wuji mcp status filesystem
wuji mcp remove filesystem
```

Transport modes:
- `--stdio` — foreground proxy; use in OpenCode as `command: ["wuji", "mcp", "start", "<name>", "--stdio"]`
- `--http` / `--sse` — background gateway via supergateway; OpenCode connects with `type: remote` + URL
- default (no flag) — background subprocess with logs only (no MCP client attached)

Logs for background modes: `.wuji/runtime/mcp/<name>.log`.
Remote servers (`url`) with `--stdio` use `mcp-remote`; with `--http` only health-check the URL.

```yaml
# .wuji/config.yaml (optional — created by `wuji config set`)
default_driver: dummy
capability_drivers:
  text: vllm
  image: dummy
drivers:
  resources:
    total_ram_mb: 51200        # system RAM budget (auto-detected if omitted)
    total_vram_mb: 12288       # GPU VRAM budget (optional)
    lazy_restore: true         # default: don't reload text model after image jobs
    evict_and_restore: true
    unload_transient_models: true
    model_ram_mb:
      my-chat.gguf: 14000
    model_vram_mb:
      my-chat.gguf: 6000
      sd-xl: 8000
  vllm:
    default_model: meta-llama/Llama-3.2-1B
    host: 127.0.0.1
    port: 8000
    startup_timeout_seconds: 300
  llama:
    default_model: my-model.gguf
    models_dir: ../driver/llama/models
    inference_port: 8080
    ollama_api: http://127.0.0.1:11434
```

## Resource management (automatic)

**No manual setup required.** Every generate command (`wuji text`, `wuji image`, `wuji video`, `wuji audio`, `wuji voice`) runs through the resource scheduler when RAM or VRAM can be detected. If `drivers.resources` is missing from config, Wuji auto-detects system RAM (`/proc/meminfo`) and GPU VRAM (`nvidia-smi`) and applies conservative budgets (80% RAM, 90% VRAM).

### What happens on each command

1. **Lock** — file lock at `.wuji/runtime/scheduler.lock` coordinates parallel `wuji` processes.
2. **Plan** — if the requested model does not fit the budget, conflicting loaded models are evicted (unloaded) first.
3. **Run** — the command executes.
4. **Cleanup** — transient capabilities (image, video, audio, mesh) unload their model after the job.
5. **Restore** — with `lazy_restore: true` (default), evicted models stay unloaded until the next matching command (e.g. `wuji text` after `wuji image`). No immediate reload spike.

### Example: chat + image without OOM

```
Terminal A: wuji text -p "hello"      → loads chat model (vllm/llama)
Terminal B: wuji image -p "a cat"     → evicts chat, runs image, frees image RAM
Terminal A: wuji text -p "again"      → reloads chat on demand
```

State persists in `.wuji/runtime/scheduler.json` across CLI invocations.

### CLI

```bash
wuji resources explain   # full walkthrough
wuji resources status    # budgets, loaded slots, restore queue
wuji resources load      # manually preload a model (vllm, llama, …)
wuji resources unload    # manually free memory (vllm | llama | all)
wuji resources detect    # show detected RAM/VRAM
wuji resources init      # persist detected budgets to config (optional)
wuji resources kill      # alias for unload all
```

`wuji resources init` only pins budgets in config — scheduling stays automatic either way.

## Proto / codegen

```bash
make install-tools  # buf, protoc plugins
make proto          # regenerates api/proto/v1/*.pb.go
```

## Common files when changing…

| Task | Start here |
|------|------------|
| New CLI flag / command | `internal/cli/` |
| MCP server management | `internal/mcp/`, `internal/config/mcp.go`, `wuji mcp` |
| New request field | `internal/driver/types.go`, `api/proto/v1/driver.proto`, `*_proto.go`, then `make proto` |
| New capability | `internal/capability/`, interfaces in `internal/driver/capabilities.go`, proto RPC, host + client in `internal/driver/grpc/` |
| Resource scheduling / OOM | `internal/scheduler/`, `internal/sysmem/`, `internal/core/resources.go` |

## Module

`github.com/coditary/wuji` — Go 1.23, Cobra CLI, gRPC/protobuf.
