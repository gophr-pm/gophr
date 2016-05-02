#!/bin/bash

perl -pi -e "s/%%ip%%/$(hostname -I)/g" /etc/cassandra/cassandra.yaml
export CLASSPATH=/kubernetes-cassandra.jar
cassandra -f
