#!/usr/bin/env sh

url="$1"
path="$2"
date="$(date '+%Y%m%d-%H%M%S')"

echo "${url} ${path} ${date}"

mkdir -p "${path}"

timeout 30 ffmpeg -rtsp_transport tcp \
    -y -i "${url}" \
    -f image2 -vf fps=fps=2 -s 1024x768 "${path}/cap-${date}.jpg"

exit 0
