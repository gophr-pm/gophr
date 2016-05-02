FROM golang:alpine
# Install git
RUN set -ex && apk add --update 'git'
# Specify the volume mount point as current dir
VOLUME ["/go/src/github.com/skeswa/gophr"]
WORKDIR /go/src/github.com/skeswa/gophr
# Get all the appropriate dependencies
CMD cd ./indexer && go get -d -v && go run *.go
