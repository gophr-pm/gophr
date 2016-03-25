FROM golang:alpine
# Install gin
RUN set -ex && apk add --update 'git'
# Install git
RUN go get github.com/codegangsta/gin
# Set the environment variables for dev
ENV PORT="3000"
# Specify the exposed port
EXPOSE 3000
# Specify the volume mount point as current dir
VOLUME ["/go/src/app"]
WORKDIR /go/src/app
# Get all the appropriate dependencies
CMD go get -d -v && gin
