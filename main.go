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

var scanWorkers = flag.Int("w", 1024, "workers count")
var generateIps = flag.Int64("gw", -1, "generate N random WAN IPs")
var resultCallback = flag.String("callback", "", "callback to run as shell command. Use {result} as template")
var resultCallbackTimeout = flag.Int("ct", 30, "callback timeout in seconds")

var scanRtsp = flag.Bool("rtsp", false, "check rtsp")
var rtspFuzzDict = flag.String("rtspd", "./data/rtsp-paths.txt", "RTSP dictionary to fuzz")

var scanPorts = flag.String("ports", "", "scan ports on every rarget")
var portScanTimeout = flag.Int("pst", 700, "portscan timeout in milliseconds")

func main() {
	flag.Parse()

	ipGenerator := generators.NewIPGenerator(128, *generateIps)
	processor := services.NewProcessor(ipGenerator, *scanWorkers)

	callbackTimeout := time.Second * time.Duration(*resultCallbackTimeout)

	if *scanRtsp {
		os.Stderr.WriteString(fmt.Sprintln("[i] Using rtsp"))
		paths, err := loaders.FileToArray(*rtspFuzzDict)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintln("[E]", err))
			os.Exit(1)
		}
		rtspService := services.NewRTSPService(554, paths)
		processor.AddService(rtspService)
	}

    if *scanPorts != "" {
        ports := utils.ParseRange(*scanPorts)
        processor.AddService(services.NewPortscanService(ports, time.Millisecond * time.Duration(*portScanTimeout)))
    }

	if len(processor.Services()) == 0 {
		processor.AddService(services.NewDummyService())
	}

	wg := new(sync.WaitGroup)

	for result := range processor.Process() {
		fmt.Println(&result)
		if *resultCallback != "" {
			wg.Add(1)
			cmd := result.ReplaceVars(*resultCallback)
			go utils.RunCommand(cmd, wg, callbackTimeout)
		}
	}

	wg.Wait()
}
