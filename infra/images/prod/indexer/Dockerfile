FROM golang:alpine
# Install gin
RUN set -ex && apk add --update 'git'

COPY . /go/src/github.com/skeswa/gophr
WORKDIR /go/src/github.com/skeswa/gophr/indexer
RUN echo "* * * * * printenv && /go/src/github.com/skeswa/gophr/indexer/indexer -e prod -p 3000 -m /go/src/github.com/skeswa/gophr/common/db/migrations >> /var/log/indexer" > /var/spool/cron/crontabs/root
RUN go get -d -v && go build .

# Specify the exposed port
EXPOSE 3000

# Get all the appropriate dependencies
CMD crond -l 2 -f
