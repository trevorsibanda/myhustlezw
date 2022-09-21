#!/bin/bash

source ./prod.env
export RABBIT_MQ_VOLUME_DIR="$(pwd)/$RABBIT_VOLUME_DIR"
echo "Starting rabbitmq with persistance dir: $RABBITMQ_MQ_VOLUME_DIR"
docker run --rm -d --hostname $RABBITMQ_HOST -p 15672:15672 -p 5671:5671 -p 5672:5672 -p 15671:15671 -p 25672:25672  --name myhustle-rabbit -e RABBITMQ_DEFAULT_USER="$RABBITMQ_USERNAME" -e RABBITMQ_DEFAULT_PASS="$RABBITMQ_PASSWORD" rabbitmq:3-management
