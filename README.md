# Netstalking things in GO

- RTSP fuzzer

## Build

```sh
go build
```

## Usage

```sh
./go-ns -rtsp -w 4096 -callback './callbacks/capture.sh "{result}" "/sdcard/Pictures/RTSP/" "{slug}"'
```
