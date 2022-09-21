#!/usr/bin/bash

input=$1
output_dir=$2
echo "Usage:\n ./encode.sh input output" | cat >> log.txt
echo "\nInput: $input" | cat >> log.txt
echo "Output dir: $output_dir" | cat >> log.txt 
cd $output_dir
ffmpeg -y -i $input -c:a aac -b:a 128k -muxdelay 0 -f segment -sc_threshold 0 -segment_time 7 -segment_list "file.m3u8" -segment_format mpegts "file%d.m4a"