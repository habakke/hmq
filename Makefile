BINARY          := hmq
ROOT_DIR        := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR       := $(ROOT_DIR)/build
BUILD_DIR       := $(ROOT_DIR)/dist
VERSION         := $(shell cat VERSION)
GITSHA          := $(shell git rev-parse --short HEAD)

.PHONY: build clean start lint staticcheck test fmt release-test release

prepare:
	mkdir -p $(BUILD_DIR)

lint:
	go get -u golang.org/x/lint/golint
	$(shell go list -f {{.Target}} golang.org/x/lint/golint) ./...

check:
	go get -u honnef.co/go/tools/cmd/staticcheck
	$(shell go list -f {{.Target}} honnef.co/go/tools/cmd/staticcheck) ./...

test: prepare
	go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

build: prepare
	goreleaser build --snapshot --rm-dist

start:
	go run $(ROOT_DIR)/cmd/$(BINARY)/main.go

profile:
	go tool pprof -http=:7777 cpuprofile

clean:
	rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...

release-test: export GITHUB_SHA=$(GITSHA)
release-test:
	goreleaser release --skip-publish --snapshot --rm-dist

release: export GITHUB_SHA=$(GITSHA)
release: release-test
	git tag -a $(VERSION) -m "Release" && git push origin $(VERSION)
