FROM alpine:3.7
RUN apk --no-cache add ca-certificates docker
RUN update-ca-certificates
WORKDIR /buildctl/
COPY bin/buildctl .
COPY conf/buildctl.yml .
ENTRYPOINT ["/buildctl/buildctl"]