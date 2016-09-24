#!/bin/bash

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
