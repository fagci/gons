#!/usr/bin/env sh

url="$1"
dir="$2"
slug="$3"
path="${dir}/${slug}.png"

mkdir -p "${dir}"

timeout 20 ffmpeg -rtsp_transport tcp \
    -i "$url" -vf fps=1/3 -pix_fmt yuvj420p -nostdin "$path" -y -loglevel error 2>&1
