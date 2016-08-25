
# gophr
[![Go Report Card](https://goreportcard.com/badge/github.com/skeswa/gophr)](https://goreportcard.com/report/github.com/skeswa/gophr)
[![GoDoc](https://godoc.org/github.com/skeswa/gophr/common?status.svg)](https://godoc.org/github.com/skeswa/gophr/common)

An end-to-end package management solution for Go.

## Setting up the dev environment
Firstly, make sure you have `docker` (and `docker-machine`) installed. You should be able to run this command:
```sh
$ docker && docker-machine || echo "Not properly installed"
```
If not, resources on docker setup can be found here:
- [Mac OSX](https://docs.docker.com/engine/installation/mac/)
- [Ubuntu](https://docs.docker.com/engine/installation/linux/ubuntulinux/)
- [Windows](https://docs.docker.com/engine/installation/windows/)
- [Other](https://docs.docker.com/engine/installation/)

Then, since the `gophrctl` tool needs to be installed for the environment to function, run the following script:
```sh
$ cd $GOPHR_REPO/infra/bin/setup/ && ./install
```
Next, the dev environment needs initialization. So, run the following:
```sh
$ gophrctl init
```
Afterwards, you can run the following for information on how to use `gophrctl`:
```sh
$ gophrctl --help
```
## Running the dev server
Before the server can start, the dev docker images needs to be built. So, run the following (keep in mind that it will take a while):
```sh
$ gophrctl build
```
After the images are built, you can run:
```sh
$ gophrctl up
```
This starts every component of the dev server. Lastly, in order to compile the web application, run:
```sh
$ gophrctrl web
```
At this point, you should be able to open https://gophr.dev/ in your favorite browser.

## Running the indexer

We've already created a new default db keyspace `gophr` with a table named `packages`in our cassandra db via our dockerfile. Now we need to load data in to work with. The indexer pulls go package metadata from various sources and compiles them into packageDTOs and inserts them into the DB.

Navigate to the `indexer` folder:
```sh
$ cd $GOPHR_REPO/indexer
```

Build the indexer binary:
```sh
$ go build && ./indexer
```

## Using Gophr-cli

Now that you've completely setup your `gophr dev environment` and have loaded data via the `indexer` you're ready to use the [Gophr-cli](https://github.com/Shikkic/gophr-cli) tool that integrates with the Gophr dev environment.

