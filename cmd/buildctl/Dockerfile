FROM golang
ARG version
COPY pkg/ /go/src/github.com/mobingilabs/pullr/pkg/
COPY vendor/ /go/src/github.com/mobingilabs/pullr/vendor/
COPY cmd/buildctl/ /go/src/github.com/mobingilabs/pullr/cmd/buildctl/
WORKDIR /go/src/github.com/mobingilabs/pullr/cmd/buildctl/
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o buildctl .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /buildctl/
COPY --from=0 /go/src/github.com/mobingilabs/pullr/cmd/buildctl/buildctl .
ENTRYPOINT ["/buildctl/buildctl"]