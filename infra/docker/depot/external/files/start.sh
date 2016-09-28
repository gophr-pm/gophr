#!/bin/bash

echo "Starting fcgi..."
/etc/init.d/spawn-fcgi start
echo "Starting nginx..."
nginx -g 'daemon off;'
