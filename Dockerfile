# syntax = docker/dockerfile:1-experimental

FROM --platform=$BUILDPLATFORM golang:1.16-alpine AS base
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

FROM scratch AS unit-test-coverage
COPY --from=unit-test /src/build/cover.out /cover.out

FROM busybox:musl AS bin
COPY --from=build /src/build/hmq .
EXPOSE 1883
EXPOSE 8080
CMD ["/hmq"]
