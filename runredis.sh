#!/bin/bash

source ./prod.env
export CACHE_REDIS_VOLUME_DIR="$(pwd)/$CACHE_REDIS_VOLUME_DIR"
echo "Starting redis server with persistance dir: $CACHE_REDIS_VOLUME_DIR"
docker run --name myhustle-redis -v $CACHE_REDIS_VOLUME_DIR:/data -p 6379:6379 --rm -d redis redis-server --save 60 1 --loglevel warning
