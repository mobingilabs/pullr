FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /eventsrv/
COPY bin/eventsrv .
COPY conf/eventsrv.yml .
ENTRYPOINT ["/eventsrv/eventsrv"]
