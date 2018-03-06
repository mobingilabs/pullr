FROM scratch
WORKDIR /app
ADD certs/ca-certificates.crt /etc/ssl/certs
ADD conf/trackersvc.yml .
ADD bin/trackersvc .
ENTRYPOINT ["/app/trackersvc"]
CMD ["listen"]
