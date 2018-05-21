FROM golang:1.9
WORKDIR /go/src/github.com/mobingilabs/pullr
COPY cmd ./cmd
COPY vendor ./vendor
COPY pkg ./pkg
COPY conf ./conf
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/apisrv ./cmd/apisrv

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /apisrv/
COPY --from=0  /go/src/github.com/mobingilabs/pullr/bin/apisrv .
COPY conf/pullr.yml .
ENTRYPOINT ["/apisrv/apisrv"]
