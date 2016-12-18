FROM golang:alpine

# Set the default timezone to EST.
ENV TZ=America/New_York
RUN set -ex \
  && apk add --update tzdata \
  && cp "/usr/share/zoneinfo/$TZ" /etc/localtime \
  && echo $TZ > /etc/timezone

# We need git for pulling dependencies and bash for debugging.
RUN set -ex && apk add --update 'git' && apk add 'bash'

# Copy in requisite files.
COPY ./scheduler/master /go/src/github.com/gophr-pm/gophr/scheduler/master
COPY ./infra /go/src/github.com/gophr-pm/gophr/infra
COPY ./lib /go/src/github.com/gophr-pm/gophr/lib

# Build source and move things around.
RUN cd /go/src/github.com/gophr-pm/gophr/scheduler/master \
    && echo -e "\nFetching dependencies...\n" \
    && go get -d -v \
    && echo -e "\nBuilding the binary...\n" \
    && go build -v -o gophr-scheduler-master-binary \
    && chmod +x ./gophr-scheduler-master-binary \
    && echo -e "\nMoving things around...\n" \
    && mkdir /gophr \
    && mv ./gophr-scheduler-master-binary /gophr/gophr-scheduler-master-binary \
    && cd /gophr \
    && rm -rf /go

# Since its no longer necessary, remove git.
RUN set -ex && apk del git

# Set the environment variables
ENV PORT="3000"
ENV GOPHR_ENV="prod"
ENV GOPHR_SECRETS_PATH="/secrets"

VOLUME ["/secrets"]

# Set the generated directory as the work directory.
WORKDIR /gophr

CMD ./gophr-scheduler-master-binary --port "$PORT"
