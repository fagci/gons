package main

import (
	"errors"
	"fmt"
	"go_ns/src/gen"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const PORT = "554"

const RTSP_OPT = "OPTIONS * RTSP/1.0\r\nCSeq: 1\r\n\r\n"

func rtsp_req(c net.Conn, req string) (int, error) {
	if _, e := c.Write([]byte(RTSP_OPT)); e != nil {
		return 0, e
	}

	m := make([]byte, 1024)
	if _, e := c.Read(m); e != nil {
		return 0, e
	}

	f := strings.Fields(string(m))
	if len(f) > 2 {
		return strconv.Atoi(f[1])
	}
	return 0, errors.New("Bad response")
}

func check(paths []string, ip string, port string) {
	addr := ip + ":" + port

	d := net.Dialer{Timeout: time.Second * 3}

	c, e := d.Dial("tcp", addr)

	if e != nil {
		return
	}

	defer c.Close()

	var code int
	code, e = rtsp_req(c, RTSP_OPT)
	if e != nil || code != 200 {
		return
	}

	req := "DESCRIBE / RTSP/1.0\r\nCSeq: 2\r\n\r\n"

	code, e = rtsp_req(c, req)
	if e != nil || code == 401 {
		return
	}

	if code == 200 {
		fmt.Printf("rtsp://%s/\n", addr)
		return
	}

	for i, path := range paths {
		req := fmt.Sprintf("DESCRIBE %s RTSP/1.0\r\nCSeq: %d\r\n\r\n", path, i+3)
		log.Println(req)
		code, e = rtsp_req(c, req)
		if e != nil || code == 401 {
			return
		}
		if code == 200 {
			fmt.Printf("rtsp://%s%s\n", addr, path)
			return
		}
	}
}

func worker(paths []string) {
	gen_ip := gen.WanIpGenerator()
	for {
		check(paths, gen_ip().String(), PORT)
	}
}

func main() {
	paths := []string{
		"/1",
		"/0/1:1/main",
		"/live/h264",
		"/live",
		"/h264/ch1/sub/av_stream",
		"/stream1",
		"/live.sdp",
		"/image.mpg",
		"/axis-media/media.amp",
		"/1/stream1",
		"/ch01.264",
		"/live1.sdp",
		"/stream.sdp",
	}

	for i := 0; i < 1500; i++ {
		go worker(paths)
	}

	runtime.Goexit()
}
