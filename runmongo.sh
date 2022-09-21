#!/bin/bash

source ./prod.env
export MONGODB_VOLUME_DIR="$(pwd)/$MONGODB_VOLUME_DIR"
echo "Starting mongodb with volume dir $MONGODB_VOLUME_DIR"
docker run -d  --name myhustle-mongo \
    -e MONGO_INITDB_ROOT_USERNAME=$MONGODB_USERNAME \
    -e MONGO_INITDB_ROOT_PASSWORD=$MONGODB_PASSWORD \
    -p 27017:27017 \
    --rm \
    -v $MONGODB_VOLUME_DIR:/data/db \
    mongo
