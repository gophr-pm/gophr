FROM golang:alpine
# Install sparkplug
RUN set -ex && apk add --update 'git'
RUN go get github.com/skeswa/sparkplug
# Set the environment variables for dev
ENV PORT="3000"
ENV GOPHR_ENV="dev"
ENV GOPHR_MIGRATIONS_PATH="/go/src/github.com/skeswa/gophr/common/db/migrations"
# Specify the exposed port
EXPOSE 3000
# Specify the volume mount point as current dir
VOLUME ["/go/src/github.com/skeswa/gophr"]
WORKDIR /go/src/github.com/skeswa/gophr
# Get all the appropriate dependencies
CMD cd ./indexer \
    && go get -d -v \
    && sparkplug
