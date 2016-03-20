FROM golang:alpine
# Install gin
RUN set -ex && apk add --update 'git'
# Install git
RUN go get github.com/codegangsta/gin
# Set the environment variables for dev
ENV PORT="80"
# Specify the exposed port
EXPOSE 80
# Specify the volume mount point as current dir
VOLUME ["/go/src/app"]
WORKDIR /go/src/app
# Get all the appropriate dependencies
CMD go get -d -v && gin
