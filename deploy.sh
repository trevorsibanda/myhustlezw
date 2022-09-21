#!/bin/bash


DIST="$(pwd)/dist/$(date '+%a_%m_%y_%T')"
echo "Deploying build to $DIST"
mkdir $DIST -p
echo "Building myhustle app"
export GOOS=linux GOARCH=amd64
go build
cp myhustlezw $DIST/
echo "Copying static site"
mkdir -p $DIST/static/website/
cp static/website $DIST/static/ -R
echo "Copying built dashboard"
mkdir -p $DIST/static/dashboard/dist
cp static/dashboard/build  $DIST/static/dashboard/ -R
echo "Copying support scripts"
cp scripts $DIST -R
echo "Copying build scripts"
cp *.sh $DIST/
echo "Copying prod env"
cp prod.env $DIST
echo "Copying data dir"
mkdir -p $DIST/data/redis $DIST/data/mongodb $DIST/data/rabbit $DIST/data/uploads 
echo "Preparing zipped deployment"
cd $DIST
zip -r ~/myhustle.zip .



