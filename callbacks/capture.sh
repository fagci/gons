#!/usr/bin/env sh

url="$1"
dir="$2"
slug="$3"
path="${dir}/${slug}.jpg"

mkdir -p "${dir}"

timeout 30 \
    ffmpeg -rtsp_transport tcp \
    -i "$url" \
    -frames:v 1 -an \
    -y "$path" \
    -loglevel warning 2>&1
