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
    -loglevel warning 2>&1 \
    && ( \
      (hash exiftool && exiftool -Comment="$url" "$path" && rm "${path}_original") \
      || \
      (hash jhead && jhead -cl "$url" "$path") \
    )
