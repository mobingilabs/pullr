FROM golang
ARG version
COPY pkg/ /go/src/github.com/mobingilabs/pullr/pkg/
COPY vendor/ /go/src/github.com/mobingilabs/pullr/vendor/
COPY cmd/eventsrv/ /go/src/github.com/mobingilabs/pullr/cmd/eventsrv/
WORKDIR /go/src/github.com/mobingilabs/pullr/cmd/eventsrv/
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-X github.com/mobingilabs/pullr/cmd/eventsrv/app.version=$version" -a -installsuffix cgo -o eventsrv .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /eventsrv/
COPY --from=0 /go/src/github.com/mobingilabs/pullr/cmd/eventsrv/eventsrv .
ENTRYPOINT ["/eventsrv/eventsrv"]
