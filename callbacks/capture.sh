#!/usr/bin/env sh

url="$1"
dir="$2"
slug="$3"
path="${dir}/${slug}.png"

mkdir -p "${dir}"

timeout 15 ffmpeg -loglevel error \
        -y -rtsp_transport tcp -i "${url}" \
        -pix_fmt yuvj422p -an -t 1 -r 1 \
        "${path}" 2>&1

