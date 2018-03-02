FROM alpine:3.7
RUN apk --no-cache add ca-certificates docker
RUN update-ca-certificates
WORKDIR /trackersvc/
COPY bin/trackersvc .
COPY conf/trackersvc.yml .
ENTRYPOINT ["/trackersvc/trackersvc"]
