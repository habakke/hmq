CMD=hmq
BINARY=hmq
IMAGE=hmq
VERSION=latest
ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build
all: build

build-all: clean build-amd64 build-arm build-arm64

build-amd64:
		mkdir -p $(BUILD_DIR)
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .
		docker build -t habakke/$(IMAGE):amd64-$(VERSION) .
		docker push habakke/$(IMAGE):amd64-$(VERSION)

build-arm:
		mkdir -p $(BUILD_DIR)
		GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .
		docker build -t habakke/$(IMAGE):arm-$(VERSION) .
		docker push habakke/$(IMAGE):arm-$(VERSION)

build-arm64:
		mkdir -p $(BUILD_DIR)
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .
		docker build -t habakke/$(IMAGE):arm64-$(VERSION) .
		docker push habakke/$(IMAGE):arm64-$(VERSION)

start:
	go run $(ROOT_DIR)/main.go

clean:
	rm -rf $(BUILD_DIR)

manifest:
		# Create and push the multi-arch manifest
		docker manifest create habakke/$(IMAGE):$(VERSION) habakke/$(IMAGE):amd64-$(VERSION) habakke/$(IMAGE):arm-$(VERSION) habakke/$(IMAGE):arm64-$(VERSION)
		docker manifest annotate habakke/$(IMAGE):$(VERSION) habakke/$(IMAGE):amd64-$(VERSION) --arch amd64 --os linux
		docker manifest annotate habakke/$(IMAGE):$(VERSION) habakke/$(IMAGE):arm-$(VERSION)  --arch arm --os linux
		docker manifest annotate habakke/$(IMAGE):$(VERSION) habakke/$(IMAGE):arm64-$(VERSION) --arch arm64 --os linux
		docker manifest push --purge habakke/$(IMAGE):$(VERSION)
