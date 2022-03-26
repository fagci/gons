package main

import (
	"flag"
	"fmt"
	"go-ns/src/generators"
	"go-ns/src/loaders"
	"go-ns/src/services"
	"go-ns/src/utils"
	"os"
	"sync"
)

var scanRtsp = flag.Bool("rtsp", false, "check rtsp")
var fuzzDict = flag.String("d", "./data/rtsp-paths.txt", "dictionary to fuzz")
var scanWorkers = flag.Int("w", 1024, "workers count")
var resultCallback = flag.String("callback", "", "callback to run as shell command. Use {result} as template")

func main() {
	flag.Parse()

	paths, err := loaders.FileToArray(*fuzzDict)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	ipGenerator := generators.NewIPGenerator(128)
	processor := services.NewProcessor(ipGenerator, *scanWorkers)

	if *scanRtsp {
		fmt.Println("using rtsp")
		rtspService := services.NewRTSPService(554, paths)
		processor.AddService(rtspService)
	}

	if len(processor.Services()) == 0 {
		fmt.Println("Specify at least one service to check")
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)

	for result := range processor.Process() {
		fmt.Println(result.Url.String())
		if *resultCallback != "" {
			wg.Add(1)
			cmd := result.ReplaceVars(*resultCallback)
			go utils.RunCommand(cmd, wg)
		}
	}
	wg.Wait()
}
