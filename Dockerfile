FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Install Dependencies
RUN apk add --no-cache make git

WORKDIR /src
COPY * .
RUN echo $(ls -la)
RUN make build

FROM busybox:musl
COPY --from=build /build/hmq .
EXPOSE 1883

CMD ["/hmq"]
