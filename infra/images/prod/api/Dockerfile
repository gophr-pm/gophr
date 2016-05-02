FROM golang:alpine
# Install git
RUN set -ex && apk add --update 'git'
# Set the environment variables for dev
ENV GOPHR_PORT="3000"
ENV GOPHR_ENV="staging"
ENV GOPHR_DOMAIN="gophr.pm"
# Specify the exposed port
EXPOSE 3000
# Copy the source code
COPY ../../../../ /go/src/github.com/skeswa/gophr
WORKDIR /go/src/github.com/skeswa/gophr/api
RUN go get -d -v && go build -o server
# Get all the appropriate dependencies
CMD ./server
