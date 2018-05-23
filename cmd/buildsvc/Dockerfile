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
ADD bin/buildsvc .
ADD conf/pullr.yml .
ENTRYPOINT ["/buildsvc/buildsvc"]