package main

import (
	"fmt"
	"go_ns/src/gen"
	"go_ns/src/loaders"
	"go_ns/src/svc"
	"os"
	"runtime"
)

func worker(generator *gen.IPGenerator, paths []string) {
	for ip := range generator.GenerateWAN() {
		rtsp := svc.NewRTSP(ip.String() + ":554")
		for path := range rtsp.CheckPaths(paths) {
			fmt.Println(path)
		}
	}
}

func main() {
	generator := gen.NewIPGenerator(1024)
	loader := loaders.NewDictLoader()
	paths, err := loader.Load("./data/rtsp-paths.txt") // TODO: make possible to use different path
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	for i := 0; i < 1024; i++ {
		go worker(generator, paths)
	}

	runtime.Goexit()
}
