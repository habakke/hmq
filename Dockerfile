FROM scratch
WORKDIR /
COPY /build/hmq .
EXPOSE 1883

CMD ["/hmq"]
