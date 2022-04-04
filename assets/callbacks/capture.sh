#!/usr/bin/env bash

url="$1"
dir="$2"
slug="$3"
file_path="${dir}/${slug}.jpg"

mkdir -p "${dir}"

function capture() {
    local url="$1"
    local path="$2"
    ffmpeg -rtsp_transport tcp -i "$url" -frames:v 1 -an -y "$path" -loglevel warning
}

function add_exif_comment() {
    local file="$1"
    local comment="$2"
    if hash exiftool; then
        exiftool -q -Comment="$url" "$file_path" && rm "${file_path}_original" > /dev/null
    elif hash jhead; then
        jhead -cl "$url" "$file_path" > /dev/null
    fi
}

capture "$url" "$file_path" \
    && add_exif_comment "$file_path" "$url" \
    && echo "[+] $url"

