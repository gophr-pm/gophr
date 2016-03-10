FROM golang:1.6-onbuild
EXPOSE 3000
ENTRYPOINT ["go", "run", "*.go"]
