FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /apisrv/
ADD bin/apisrv .
ADD conf/apisrv.yml .
ENTRYPOINT ["/apisrv/apisrv"]