#!/usr/bin/env bash

set -eu

rm -rf rootfs || true
docker plugin disable -f s3vol || true
docker plugin rm -f s3vol || true

docker build --label driver="s3vol" -t rootfsimage .
id=$(docker create rootfsimage true)
mkdir -p rootfs
docker export "$id" | tar -x -C rootfs
docker rm -vf "$id"
docker rmi rootfsimage

docker plugin create s3vol .
rm -rf rootfs || true
