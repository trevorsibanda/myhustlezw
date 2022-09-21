#!/bin/bash

source prod.env
echo using key $IMGPROXY_KEY $IMGPROXY_SALT
sudo docker run -p 8089:8080 \
     -v $MYHUSTLE_LOCALFILES_DIR:/tmp \
     -e IMGPROXY_LOCAL_FILESYSTEM_ROOT=/tmp/ \
     -e IMGPROXY_SALT="$IMGPROXY_SALT" \
     -e IMGPROXY_KEY="$IMGPROXY_KEY" \
     -e IMGPROXY_S3_ENDPOINt=$S3_ENDPOINT \
     -e IMGPROXY_USE_S3=true \
     -e AWS_REGION=$S3_REGION \
     -e AWS_ACCESS_KEY=$S3_KEY \
     -e AWS_SECRET_ACCESS_KEY=$S3_SECRET darthsim/imgproxy 

