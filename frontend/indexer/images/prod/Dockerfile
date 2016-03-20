FROM golang:alpine
# Install gin
RUN ["go", "get", "github.com/codegangsta/gin"]
# Set the environment variables for dev
ENV PORT="80"
# Specify the volume mount point as current dir
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
# Run the API server
ENTRYPOINT ["gin"]
