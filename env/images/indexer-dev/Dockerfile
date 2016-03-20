FROM golang:alpine

# TODO(Sandile): cron stuff

# Just run it for now
# Install gin
RUN set -ex && apk add --update 'git'
# Specify the volume mount point as current dir
VOLUME ["/go/src/app"]
WORKDIR /go/src/app
# Get all the appropriate dependencies
CMD go get -d -v && go run *.go
