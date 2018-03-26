FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /buildsvc/
ADD bin/buildsvc .
ADD conf/pullr.yml .
ENTRYPOINT ["/buildsvc/buildsvc"]
