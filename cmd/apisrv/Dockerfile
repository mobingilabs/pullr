FROM golang
ARG version
COPY pkg/ /go/src/github.com/mobingilabs/pullr/pkg/
COPY vendor/ /go/src/github.com/mobingilabs/pullr/vendor/
COPY cmd/apisrv/ /go/src/github.com/mobingilabs/pullr/cmd/apisrv/
WORKDIR /go/src/github.com/mobingilabs/pullr/cmd/apisrv/
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-X github.com/mobingilabs/pullr/cmd/apisrv/app.version=$version" -a -installsuffix cgo -o apisrv .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /apisrv/
COPY --from=0 /go/src/github.com/mobingilabs/pullr/cmd/apisrv/apisrv .
ENTRYPOINT ["/apisrv/apisrv"]