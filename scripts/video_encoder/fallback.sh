#!/bin/bash

input=$1
output=$2
echo "Usage:\n ./runner.sh input output" | cat >> log.txt
echo "\nInput: $input" | cat >> log.txt
echo "Output: $output" | cat >> log.txt 
ffmpeg -i $input -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls $output/file.m3u8
