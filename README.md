# gophr
An end-to-end package management solution for Go.

## Setting up the dev environment
Firstly, make sure you have `docker` (and `docker-machine`) installed. You should be able to run this command:
```
$ docker && docker-machine || echo "Not properly installed"
```
If not resources on docker setup can be found [here](https://docs.docker.com/engine/installation/)
- (Mac OSX)[https://docs.docker.com/engine/installation/mac/]
- (Windows)[https://docs.docker.com/engine/installation/windows/]
- (Ubuntu)[https://docs.docker.com/engine/installation/linux/ubuntulinux/]

Then, since the `gophrctl` tool needs to be installed for the environment to function, run the following script:
```sh
$ cd $GOPHR_REPO/infra/bin/setup/ && ./install
```

Afterwards, for instructions on how to use it, run:
```sh
$ gophrctl --help
```
