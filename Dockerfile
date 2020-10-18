FROM --platform=$BUILDPLATFORM golang:1.15-alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Install Dependencies
RUN apk add --no-cache make git

# Build executable
WORKDIR /src
COPY / .
RUN make build

FROM busybox:musl
COPY --from=build /src/build/hmq .
EXPOSE 1883

CMD ["/hmq"]
