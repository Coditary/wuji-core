# Wuji Core

Backend daemon and driver orchestration for Wuji.

## Layout

```
wuji-core/
  cmd/wuji-core/          # daemon binary
  cmd/wuji-driver-*/      # built-in driver binaries
  pkg/                    # libraries (also used by wuji-ai CLI)
  internal/server/        # gRPC server
  api/proto/              # protobuf contracts
  .wuji/                  # runtime config, sockets, logs

../driver/                # external driver plugins (llama, vllm, …)
../wuji-ai/               # CLI frontend
```

## Build

```bash
make build
make build-drivers   # optional plugin drivers → bin/
make test
```

## Run

```bash
./bin/wuji-core
```

See `AGENTS.md` for architecture details.
