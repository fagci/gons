package main

import (
	"flag"
	"fmt"
	"gons/src/generators"
	"gons/src/loaders"
	"gons/src/services"
	"gons/src/utils"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var generateIps = flag.Int64("gw", -1, "generate N random WAN IPs")
var scanWorkers = flag.Int("w", 1024, "workers count")

var resultCallback = flag.String("callback", "", "callback to run as shell command. Use {result} as template")
var resultCallbackConcurrency = flag.Int("cc", 30, "callback max concurrency")
var resultCallbackTimeout = flag.Int("ct", 30, "callback timeout in seconds")
var resultCallbackE = flag.Bool("dce", false, "disable callback errors")
var resultCallbackW = flag.Bool("dcw", false, "disable callback warnings")
var resultCallbackI = flag.Bool("dci", false, "disable callback info")

var scanRtsp = flag.Bool("rtsp", false, "check rtsp")
var rtspFuzzDict = flag.String("rtspd", "./data/rtsp-paths.txt", "RTSP dictionary to fuzz")

var scanPorts = flag.String("ports", "", "scan ports on every rarget")
var portScanTimeout = flag.Int("pst", 700, "portscan timeout in milliseconds")

var cpuprofile = flag.String("cpu", "", "profile cpu")
var memprofile = flag.String("mem", "", "profile mem")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			utils.EPrintln("[E]", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			utils.EPrintln("[E]", err)
			os.Exit(1)
		}
		defer f.Close()
		defer pprof.WriteHeapProfile(f)
	}

	ipGenerator := generators.NewIPGenerator(128, *generateIps)
	processor := services.NewProcessor(ipGenerator, *scanWorkers)

	callbackTimeout := time.Second * time.Duration(*resultCallbackTimeout)
	var cbFlags utils.Flags
	if !*resultCallbackE {
		cbFlags = cbFlags.Set(utils.ERR)
	}
	if !*resultCallbackW {
		cbFlags = cbFlags.Set(utils.WARN)
	}
	if !*resultCallbackI {
		cbFlags = cbFlags.Set(utils.INFO)
	}

	if *scanRtsp {
		os.Stderr.WriteString(fmt.Sprintln("[i] Using rtsp"))
		paths, err := loaders.FileToArray(*rtspFuzzDict)
		if err != nil {
			utils.EPrintln("[E]", err)
			os.Exit(1)
		}
		rtspService := services.NewRTSPService(554, paths)
		processor.AddService(rtspService)
	}

	if *scanPorts != "" {
		ports := utils.ParseRange(*scanPorts)
		processor.AddService(services.NewPortscanService(ports, time.Millisecond*time.Duration(*portScanTimeout)))
	}

	if len(processor.Services()) == 0 {
		processor.AddService(services.NewDummyService())
	}

	wg := new(sync.WaitGroup)
	guard := make(chan struct{}, *resultCallbackConcurrency)

	for result := range processor.Process() {
		if *resultCallback != "" {
			wg.Add(1)
			guard <- struct{}{}
			cmd := result.ReplaceVars(*resultCallback)
			go func() {
				utils.RunCommand(cmd, wg, callbackTimeout, cbFlags)
				<-guard
			}()
		} else {
			fmt.Println(&result)
		}
	}

	wg.Wait()
}
