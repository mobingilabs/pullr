FROM alpine:3.7
RUN apk --no-cache add ca-certificates docker
RUN update-ca-certificates

WORKDIR /app
ADD conf/buildctl.yml .
ADD bin/buildctl .
ENTRYPOINT ["/app/buildctl"]
CMD ["listen"]
