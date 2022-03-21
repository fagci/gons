package main

import (
	"flag"
	"fmt"
	"go_ns/src/gen"
	"go_ns/src/loaders"
	"go_ns/src/svc"
	"os"
)

func main() {
	scanRtsp := flag.Bool("rtsp", false, "check rtsp")
	fuzzDict := flag.String("d", "./data/rtsp-paths.txt", "dictionary to fuzz")
	flag.Parse()

	paths, err := loaders.FileToArray(*fuzzDict)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	processor := svc.NewProcessor(gen.NewIPGenerator(1024), 1024)

	if *scanRtsp {
		fmt.Println("using rtsp")
		processor.AddService(svc.NewRTSPService(554, paths))
	}

	for result := range processor.Process() {
		fmt.Println(result.URI)
	}
}
