#!/bin/bash

# check if path exists
if [[ -z "$(ls -l /cassandra_data/data | grep gophr)" ]]; then
    # If we got here, then the gophr schema has not been created yet
    for (( ; ; ))
    do
        sleep 10
        echo "Attempting to create the schema..."
        cqlsh -f /schema.cql && break
    done
    echo "Schema was initialized correctly!"
fi
