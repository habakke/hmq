# syntax = docker/dockerfile:1-experimental

FROM --platform=BUILDPLATFORM golang:1.15-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY / .

# Install Dependencies
RUN apk add --no-cache make git

# Build executable
RUN go mod download

FROM base AS build
RUN make build

FROM base AS unit-test
RUN mkdir /out && go test -v -coverprofile=/out/cover.out ./...

FROM golangci/golangci-lint:v1.31.0-alpine AS lint-base

FROM base AS lint
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./...

FROM scratch AS unit-test-coverage
COPY --from=unit-test /out/cover.out /cover.out

FROM busybox:musl AS bin
COPY --from=build /src/build/hmq .
EXPOSE 1883
CMD ["/hmq"]
