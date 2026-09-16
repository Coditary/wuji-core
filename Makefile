.PHONY: build proto tidy test run install-tools build-drivers fmt lint cover-check ci

BIN_DIR := bin
DRIVER_ROOT ?= ../../plugins/wuji
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
	$(MAKE) -C $(DRIVER_ROOT)/llama build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)
	$(MAKE) -C $(DRIVER_ROOT)/echo build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)
	$(MAKE) -C $(DRIVER_ROOT)/vllm build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)
	$(MAKE) -C $(DRIVER_ROOT)/a1111 build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)
	$(MAKE) -C $(DRIVER_ROOT)/ffmpeg build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)
	$(MAKE) -C $(DRIVER_ROOT)/raggo build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)
	$(MAKE) -C $(DRIVER_ROOT)/gorag build WUJI_BIN_DIR=$(CURDIR)/$(BIN_DIR)

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
