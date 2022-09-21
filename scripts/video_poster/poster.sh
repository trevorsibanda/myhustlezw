#!/usr/bin/bash

input=$1
output=$2
echo "Usage:\n ./poster.sh input output" | cat >> log.txt
echo "\nInput: $input" | cat >> log.txt
echo "Output: $output" | cat >> log.txt
if [ -f $output ]; then
        echo "File already exists. Will ignore" | cat >> log.txt
        exit 0
fi 
for (( i = 10; i >= 0; i-- )); do
    echo $i
    ffmpeg -y -i $input -vf "select=gt(scene\,0.$i)" -frames:v 1  $output  || true
    if [ -f $output ]; then
        echo "Created output file with $i scene diff" | cat >> log.txt
        exit 0
    fi
done
