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
var generateIps = flag.Int("gw", 0, "generate N random WAN IPs")
var resultCallback = flag.String("callback", "", "callback to run as shell command. Use {result} as template")

func main() {
	flag.Parse()

	paths, err := loaders.FileToArray(*fuzzDict)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln("[E]", err))
		os.Exit(1)
	}

	ipGenerator := generators.NewIPGenerator(128)

	if *generateIps != 0 {
		for i := 0; i < *generateIps; i++ {
			fmt.Println(ipGenerator.GenerateWANIP())
		}
		return
	}

	processor := services.NewProcessor(ipGenerator, *scanWorkers)

	if *scanRtsp {
		os.Stderr.WriteString(fmt.Sprintln("[i] Using rtsp"))
		rtspService := services.NewRTSPService(554, paths)
		processor.AddService(rtspService)
	}

	if len(processor.Services()) == 0 {
		os.Stderr.WriteString(fmt.Sprintln("[E] Specify at least one service to check"))
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
