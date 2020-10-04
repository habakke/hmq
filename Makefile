REPO=habakke
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
		docker rmi busybox:musl
		docker build -t $(REPO)/$(IMAGE):amd64-$(VERSION) --build-arg ARCH=linux/amd64 .
		docker push $(REPO)/$(IMAGE):amd64-$(VERSION)

build-arm:
		mkdir -p $(BUILD_DIR)
		GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .
		docker rmi busybox:musl
		docker build -t $(REPO)/$(IMAGE):arm-$(VERSION) --build-arg ARCH=linux/arm .
		docker push $(REPO)/$(IMAGE):arm-$(VERSION)

build-arm64:
		mkdir -p $(BUILD_DIR)
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY) -a -ldflags '-extldflags "-static"' .
		docker rmi busybox:musl
		docker build -t $(REPO)/$(IMAGE):arm64-$(VERSION) --build-arg ARCH=linux/arm64/v8 .
		docker push $(REPO)/$(IMAGE):arm64-$(VERSION)


start:
	go run $(ROOT_DIR)/main.go

clean:
	rm -rf $(BUILD_DIR)

manifest:
		# Create and push the multi-arch manifest
		docker manifest create $(REPO)/$(IMAGE):$(VERSION) $(REPO)/$(IMAGE):amd64-$(VERSION) $(REPO)/$(IMAGE):arm-$(VERSION) $(REPO)/$(IMAGE):arm64-$(VERSION)
		docker manifest annotate $(REPO)/$(IMAGE):$(VERSION) $(REPO)/$(IMAGE):amd64-$(VERSION) --arch amd64 --os linux
		docker manifest annotate $(REPO)/$(IMAGE):$(VERSION) $(REPO)/$(IMAGE):arm-$(VERSION)  --arch arm --os linux
		docker manifest annotate $(REPO)/$(IMAGE):$(VERSION) $(REPO)/$(IMAGE):arm64-$(VERSION) --arch arm64 --os linux
		docker manifest push --purge $(REPO)/$(IMAGE):$(VERSION)
