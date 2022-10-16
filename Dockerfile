FROM golang:1.17-alpine
COPY .  /go/src/httpserver
WORKDIR /go/src/httpserver/
RUN go build -o /bin/httpserver
ENTRYPOINT ["/bin/httpserver"]
EXPOSE 80
