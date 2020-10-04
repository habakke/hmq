ARG ARCH=
FROM --platform=${ARCH} busybox:musl

WORKDIR /
COPY /build/hmq .
EXPOSE 1883

CMD ["/hmq"]
