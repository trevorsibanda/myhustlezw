#!/bin/bash

input=$1
output=$2
echo "Usage:\n ./runner.sh input output" | cat >> log.txt
echo "\nInput: $input" | cat >> log.txt
echo "Output: $output" | cat >> log.txt
BASEDIR=$(dirname "$0")
echo "$BASEDIR"
$BASEDIR/video2hls --poster-filename=file.png --poster-width=640 --output-overwrite --hls-time 6  --hls-master-playlist file.m3u8 --output  $output $input | cat >> log.txt
