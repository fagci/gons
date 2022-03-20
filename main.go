package main

import (
	"fmt"
	"go_ns/src/gen"
	"go_ns/src/svc"
	"runtime"
)

func worker(generator *gen.IPGenerator) {
	for ip := range generator.GenerateWAN() {
		rtsp := svc.NewRTSP(ip.String() + ":554")
		for path := range rtsp.CheckPaths(svc.RTSP_PATHS) {
			fmt.Println(path)
		}
	}
}

func main() {
	generator := gen.NewIPGenerator(1024)
	for i := 0; i < 1024; i++ {
		go worker(generator)
	}

	runtime.Goexit()
}
