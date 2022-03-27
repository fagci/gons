package main

import (
	"flag"
	"fmt"
	"gons/src/generators"
	"gons/src/loaders"
	"gons/src/services"
	"gons/src/utils"
	"os"
	"sync"
	"time"
)

var scanRtsp = flag.Bool("rtsp", false, "check rtsp")
var fuzzDict = flag.String("d", "./data/rtsp-paths.txt", "dictionary to fuzz")
var scanWorkers = flag.Int("w", 1024, "workers count")
var generateIps = flag.Int64("gw", -1, "generate N random WAN IPs")
var resultCallback = flag.String("callback", "", "callback to run as shell command. Use {result} as template")
var resultCallbackTimeout = flag.Int("ct", 30, "callback timeout in seconds")

func main() {
	flag.Parse()

	callbackTimeout := time.Second * time.Duration(*resultCallbackTimeout)

	paths, err := loaders.FileToArray(*fuzzDict)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln("[E]", err))
		os.Exit(1)
	}

	ipGenerator := generators.NewIPGenerator(128, *generateIps)

	processor := services.NewProcessor(ipGenerator, *scanWorkers)

	if *scanRtsp {
		os.Stderr.WriteString(fmt.Sprintln("[i] Using rtsp"))
		rtspService := services.NewRTSPService(554, paths)
		processor.AddService(rtspService)
	}

	if len(processor.Services()) == 0 {
		processor.AddService(services.NewDummyService())
	}

	wg := new(sync.WaitGroup)

	for result := range processor.Process() {
		fmt.Println(result.String())
		if *resultCallback != "" {
			wg.Add(1)
			cmd := result.ReplaceVars(*resultCallback)
			go utils.RunCommand(cmd, wg, callbackTimeout)
		}
	}
	wg.Wait()
}
