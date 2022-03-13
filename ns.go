package main

import (
	"fmt"
	"go_ns/src/gen"
	"net"
	"runtime"
	"time"
)

const PORT = "80"

func check(ip string, port string) {
	addr := ip + ":" + port

	d := net.Dialer{Timeout: time.Second}
	c, e := d.Dial("tcp", addr)

	if e == nil {
		c.Close()
		fmt.Println(addr)
	}
}

func worker() {
	gen_ip := gen.WanIpGenerator()
	for {
		check(gen_ip().String(), PORT)
	}
}

func main() {
	for i := 0; i < 1024; i++ {
		go worker()
	}

	runtime.Goexit()
}
