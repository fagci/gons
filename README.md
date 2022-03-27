# Netstalking things in GO

## Features

### Netrandom

- RTSP fuzzer

## Build

```sh
go build
```

## Usage

Generate 5 random wan IPs:

```sh
./gons -gw 5
```

Netrandom find possible RTSP sources:

```sh
./gons -rtsp
```

Take snapshots from RTSP stream and write source URL in metadata:

```sh
./gons -rtsp -w 4096 -callback './callbacks/capture.sh "{result}" "/sdcard/Pictures/RTSP/" "{slug}"'
```
