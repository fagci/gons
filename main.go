package main

import (
	"flag"
	"fmt"
	"go-ns/src/generators"
	"go-ns/src/loaders"
	"go-ns/src/services"
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

	ipGenerator := generators.NewIPGenerator(1024)
	processor := services.NewProcessor(ipGenerator, 1024)

	if *scanRtsp {
		fmt.Println("using rtsp")
		processor.AddService(services.NewRTSPService(554, paths))
	}

	for result := range processor.Process() {
		fmt.Println(result)
	}
}
