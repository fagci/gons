# Netstalking things in GO

## Features

### Netrandom

- RTSP fuzzer
- random WAN IP generator
- random WAN IP port (range) scanner

## Build

```sh
go build
```

To get smaller binsry:

```sh
go build -ldflags="-s -w"
```

## Usage

Generate 5 random wan IPs:

```sh
./gons -n 5
```

Netrandom find possible RTSP sources:

```sh
./gons -s rtsp
```

Take snapshots from RTSP stream and write source URL in metadata:

```sh
./gons -s rtsp -w 4096 -cb 'bash ./callbacks/capture.sh "{result}" "/sdcard/Pictures/RTSP/" "{slug}"'
```

Scan 1024 random WAN IPs for open VNC ports:

```sh
./gons -n 1024 -s portscan -ports 5900-5902
```

## Testing

```sh
go test -v ./...
```
