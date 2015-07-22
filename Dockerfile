FROM gliderlabs/alpine:3.2
MAINTAINER Brian Hicks <brian@brianthicks.com>

RUN apk add --update ca-certificates bash

COPY . /go/src/github.com/asteris-llc/tellus
RUN apk add go git mercurial \
  && cd /go/src/github.com/asteris-llc/tellus \
  && export GOPATH=/go \
  && go get github.com/tools/godep \
  && /go/bin/godep restore \
  && go test -short ./... \
  && go build -o /bin/tellus ./cmd/tellus \
  && rm -rf /go \
  && apk del --purge go git mercurial

ENTRYPOINT ["/bin/tellus"]
