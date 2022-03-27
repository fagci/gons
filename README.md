# Netstalking things in GO

## Features

### Netrandom

- RTSP fuzzer

## Build

```sh
go build
```

## Usage

Simple:

```sh
./go-ns -rtsp
```

Take snapshots from RTSP stream and write source URL in metadata:

```sh
./go-ns -rtsp -w 4096 -callback './callbacks/capture.sh "{result}" "/sdcard/Pictures/RTSP/" "{slug}"'
```
