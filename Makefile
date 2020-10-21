REPO=habakke
CMD=hmq
BINARY=hmq
IMAGE=hmq
ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build

prepare:
	mkdir -p $(BUILD_DIR)

lint: export DOCKER_BUILDKIT=1
lint:
	@docker build . --target lint

test: export DOCKER_BUILDKIT=1
test: prepare
	@docker build . --target unit-test

build: export DOCKER_BUILDKIT=1
build: prepare
	@docker build . --target bin --output $(BUILD_DIR)

test-local: export CGO_ENABLED=0
test-local: prepare
	 go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

build-local: export CGO_ENABLED=0
build-local: prepare
	go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .

start-local: build-local
	go run $(ROOT_DIR)/main.go

clean:
	rm -rf $(BUILD_DIR)
