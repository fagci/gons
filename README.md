# Netstalking things in GO

## Features

### Netrandom

- RTSP fuzzer
- random WAN IP generator
- random WAN IP port (range) scanner
- callback command support for each result

[![DUMELS Diagram](https://www.dumels.com/api/v1/badge/e32e5a35-9583-4902-aeef-82011e033de1)](https://www.dumels.com/diagram/e32e5a35-9583-4902-aeef-82011e033de1)
[![Go Report Card](https://goreportcard.com/badge/github.com/fagci/gons)](https://goreportcard.com/report/github.com/fagci/gons)

## Build

```sh
go build
```

To get smaller binary:

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
./gons -s rtsp -cb 'bash ./assets/callbacks/capture.sh "{result}" "/sdcard/Pictures/RTSP/" "{slug}"'
```

Scan 1024 random WAN IPs for open VNC ports:

```sh
./gons -n 1024 -s portscan -ports 5900-5902
```

## Testing

```sh
go test -v ./...
```
