#!/bin/bash

set -e

# Array of all the paths on the filesystem to mount.
mounts=('/repos')

# problem fix with duplicate export entries, when running the server twice
# https://github.com/cpuguy83/docker-nfs-server/pull/12/files
echo "#NFS Exports" > /etc/exports

for mnt in "${mounts[@]}"; do
  src=$(echo $mnt | awk -F':' '{ print $1 }')
  mkdir -p $src
  echo "$src *(rw,sync,no_subtree_check,fsid=0,no_root_squash)" >> /etc/exports
done

exec runsvdir /etc/sv
