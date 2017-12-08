FROM golang

# Note that these env variables are visible via `docker history`, 'docker inspect`.
# Never upload to public registry; only ECR (current).
ARG version
ARG awsrgn
ARG awsid
ARG awssec
ENV AWS_REGION=$awsrgn \
    AWS_ACCESS_KEY_ID=$awsid \
    AWS_SECRET_ACCESS_KEY=$awssec
ADD . /go/src/github.com/mobingilabs/pullr/cmd/apisvr/
WORKDIR /go/src/github.com/mobingilabs/pullr/cmd/apisvr/
RUN go build -v -ldflags "-X github.com/mobingilabs/pullr/cmd/apisvr/cmd.version=$version"

ENTRYPOINT ["/go/src/github.com/mobingilabs/pullr/cmd/apisvr"]
CMD ["serve", "--logtostderr"]
