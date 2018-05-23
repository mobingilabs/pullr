FROM golang:1.9
WORKDIR /go/src/github.com/mobingilabs/pullr
COPY cmd ./cmd
COPY vendor ./vendor
COPY pkg ./pkg
COPY conf ./conf
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/buildsvc ./cmd/buildsvc

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /buildsvc/
COPY --from=0  /go/src/github.com/mobingilabs/pullr/bin/apisrv .
ADD conf/pullr.yml .
ENTRYPOINT ["/buildsvc/buildsvc"]