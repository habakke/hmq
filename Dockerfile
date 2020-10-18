FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Install Dependencies
RUN apk update && apk upgrade && \
    apk add --no-cache make

ADD ./ build
WORKDIR /build
RUN echo $(ls -la)
RUN echo $(ls -la /build)
RUN make build

FROM busybox:musl
COPY --from=build /build/hmq .
EXPOSE 1883

CMD ["/hmq"]
