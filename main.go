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

func main() {
	scanRtsp := flag.Bool("rtsp", false, "check rtsp")
	fuzzDict := flag.String("d", "./data/rtsp-paths.txt", "dictionary to fuzz")
	scanWorkers := flag.Int("w", 1024, "workers count")
	resultCallback := flag.String("callback", "", "callback to run as shell command. Use {result} as template")
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
