#!/bin/bash

if [[ $(nodetool status | grep $POD_IP) == *"UN"* ]]; then
  if [[ $DEBUG ]]; then
    echo "Up and normal!";
  fi
  exit 0;
else
  if [[ $DEBUG ]]; then
    echo "Not up :(";
  fi
  exit 1;
fi
