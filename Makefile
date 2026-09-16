.PHONY: build proto tidy test run install-tools build-drivers fmt lint cover-check ci

BIN_DIR := bin
COVER_MIN := 80
COVER_PKGS := \
	./pkg/auth \
	./pkg/batchplan \
	./pkg/capability \
	./pkg/loraweight \
	./pkg/modelformat \
	./pkg/scalefactor

build:
	go build -mod=mod -o $(BIN_DIR)/wuji-core ./cmd/wuji-core
	go build -mod=mod -o $(BIN_DIR)/wuji-driver-dummy ./cmd/wuji-driver-dummy
	go build -mod=mod -o $(BIN_DIR)/wuji-driver-local-rag ./cmd/wuji-driver-local-rag

build-drivers:
	$(MAKE) -C ../driver/llama build
	$(MAKE) -C ../driver/echo build
	$(MAKE) -C ../driver/vllm build
	$(MAKE) -C ../driver/a1111 build
	$(MAKE) -C ../driver/ffmpeg build
	$(MAKE) -C ../driver/raggo build
	$(MAKE) -C ../driver/gorag build

proto:
	buf generate

tidy:
	go mod tidy

test:
	go test ./...

fmt:
	@test -z "$$(gofmt -l .)"

lint:
	golangci-lint run ./...

cover-check:
	@chmod +x scripts/check-cover-min.sh
	@scripts/check-cover-min.sh $(COVER_MIN) $(COVER_PKGS)

ci: fmt lint build test cover-check

run:
	go run ./cmd/wuji-core

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
