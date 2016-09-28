#!/bin/bash

echo "Verifying that the depot volume is writable..."
if [[ -z $(mount | grep "on $GOPHR_DEPOT_PATH type nfs (rw") ]]; then
  echo "\"$GOPHR_DEPOT_PATH\" is not writable. Now exiting..."
  exit 1
else
  echo "\"$GOPHR_DEPOT_PATH\" is writable! Continuing..."
fi

echo "Starting fcgi..."
/etc/init.d/spawn-fcgi start
echo "Starting nginx in background..."
nginx
echo "Starting the depot API in the foreground..."
/gophr/wait-for-it.sh \
  -h "$GOPHR_DB_ADDR" \
  -p 9042 \
  -t 0 \
  -- \
  /gophr/gophr-depot-api-binary --port "$PORT"
