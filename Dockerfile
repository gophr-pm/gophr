FROM golang:1.6-onbuild
EXPOSE 80
ENTRYPOINT ["go", "run", "*.go"]
