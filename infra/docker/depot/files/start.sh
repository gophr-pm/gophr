#!/bin/bash

echo "Replacing \"{{API_ADDR}}\" in /etc/nginx/conf.d/git.conf with $API_SVC_SERVICE_HOST..."
sed -i "s/{{API_ADDR}}/${API_SVC_SERVICE_HOST}/g;" /etc/nginx/conf.d/git.conf
echo "Replacing \"{{ROUTER_ADDR}}\" in /etc/nginx/conf.d/git.conf with $ROUTER_SVC_SERVICE_HOST..."
sed -i "s/{{ROUTER_ADDR}}/${ROUTER_SVC_SERVICE_HOST}/g;" /etc/nginx/conf.d/git.conf

echo "Starting fcgi..."
/etc/init.d/spawn-fcgi start
echo "Starting nginx..."
nginx -g 'daemon off;'
