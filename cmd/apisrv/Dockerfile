FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /apisrv/
COPY bin/apisrv .
COPY conf/apisrv.yml .
ENTRYPOINT ["/apisrv/apisrv"]