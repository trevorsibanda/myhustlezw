#!/bin/bash

source prod.env

docker kill myhustle-mongo
docker kill myhustle-redis
docker kill myhustle-rabbit
docker kill myhustle-memcached
docker kill myhustle-payment

