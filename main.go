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

var generateIps = flag.Int64("n", -1, "generate N random WAN IPs")
var scanWorkers = flag.Int("w", 1024, "workers count")
var connTimeout = flag.Duration("t", 700*time.Millisecond, "scan connect timeout")
var scanPorts = flag.String("p", "", "scan ports on every rarget")

var resultCallback = flag.String("cb", "", "callback to run as shell command. Use {result} as template")
var resultCallbackTimeout = flag.Duration("cbt", 30*time.Second, "callback timeout")
var resultCallbackConcurrency = flag.Int("cbmc", 30, "callback max concurrency")
var resultCallbackE = flag.Bool("cbde", false, "disable callback errors")
var resultCallbackW = flag.Bool("cbdw", false, "disable callback warnings")
var resultCallbackI = flag.Bool("cbdi", false, "disable callback info")

var scanRtsp = flag.Bool("rtsp", false, "check rtsp")
var rtspFuzzDict = flag.String("rtspd", "./data/rtsp-paths.txt", "RTSP dictionary to fuzz")

var cpuprofile = flag.String("profcpu", "", "profile cpu")
var memprofile = flag.String("profmem", "", "profile mem")

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
		rtspService := services.NewRTSPService(554, paths, *connTimeout)
		processor.AddService(rtspService)
	}

	if *scanPorts != "" {
		ports := utils.ParseRange(*scanPorts)
		processor.AddService(services.NewPortscanService(ports, *connTimeout))
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
				utils.RunCommand(cmd, wg, *resultCallbackTimeout, cbFlags)
				<-guard
			}()
		} else {
			fmt.Println(&result)
		}
	}

	wg.Wait()
}
