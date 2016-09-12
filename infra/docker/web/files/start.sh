#!/bin/bash

echo "Replacing \"{{API_ADDR}}\" in /etc/nginx/nginx.conf with $API_SVC_SERVICE_HOST..."
sed -i "s/{{API_ADDR}}/${API_SVC_SERVICE_HOST}/g;" /etc/nginx/nginx.conf
echo "Replacing \"{{ROUTER_ADDR}}\" in /etc/nginx/nginx.conf with $ROUTER_SVC_SERVICE_HOST..."
sed -i "s/{{ROUTER_ADDR}}/${ROUTER_SVC_SERVICE_HOST}/g;" /etc/nginx/nginx.conf
echo "Replacing \"{{DEPOT_ADDR}}\" in /etc/nginx/nginx.conf with $DEPOT_SVC_SERVICE_HOST..."
sed -i "s/{{DEPOT_ADDR}}/${DEPOT_SVC_SERVICE_HOST}/g;" /etc/nginx/nginx.conf

echo "Starting nginx..."
nginx -g 'daemon off;'
