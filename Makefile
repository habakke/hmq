REPO=habakke
CMD=hmq
BINARY=hmq
IMAGE=hmq
ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build
all: build

prepare:
	mkdir -p $(BUILD_DIR)

test:
	@docker build . --target unit-test

lint:
	@docker build . --target lint

build: export DOCKER_BUILDKIT=1
build: prepare
	@docker build . --target bin --output $(BUILD_DIR)

start: export CGO_ENABLED=0
start: prepare build
	go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .
	go run $(ROOT_DIR)/main.go

clean:
	rm -rf $(BUILD_DIR)
