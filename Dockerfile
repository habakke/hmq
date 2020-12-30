# syntax = docker/dockerfile:1-experimental

FROM --platform=$BUILDPLATFORM golang:1.15-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY / .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Install Dependencies
RUN apk add --no-cache make git

FROM base AS build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build-local

FROM base AS unit-test
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make test-local

FROM golangci/golangci-lint:v1.31.0-alpine AS lint-base
FROM base AS lint
RUN --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./...

FROM scratch AS unit-test-coverage
COPY --from=unit-test /src/build/cover.out /cover.out

FROM busybox:musl AS bin
COPY --from=build /src/build/hmq .
EXPOSE 1883
CMD ["/hmq"]
