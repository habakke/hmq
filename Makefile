REPO=habakke
CMD=hmq
BINARY=hmq
IMAGE=hmq
ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build
all: build

build-all: clean build

build:
		mkdir -p $(BUILD_DIR)
		CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .

test:
	go test ./...

start:
	go run $(ROOT_DIR)/main.go

clean:
	rm -rf $(BUILD_DIR)
