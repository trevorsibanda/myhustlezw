#!/bin/bash

source ./prod.env
echo "Starting memcached server "
docker run --name myhustle-memcached -p 11211:11211 --rm -d memcached memcached -m 128