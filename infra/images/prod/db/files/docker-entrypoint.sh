#! /bin/bash

/run-cassandra.sh && /create-schema.sh && pkill -f cassandra
exec "$@"
