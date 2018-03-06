FROM scratch
WORKDIR /app
COPY certs/ca-certificates.crt .
COPY conf/apisrv.yml .
COPY bin/apisrv .
ENTRYPOINT ["/app/apisrv"]
CMD ["serve"]
