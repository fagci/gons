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

var randomIPsCount int64
var scanWorkers, resultCallbackConcurrency int
var connTimeout, resultCallbackTimeout time.Duration

var resultCallback string
var resultCallbackI, resultCallbackE, resultCallbackW bool

var scanPorts string

var scanRtsp bool
var rtspFuzzDict string

var cpuprofile, memprofile string

func init() {
	flag.Int64Var(&randomIPsCount, "n", -1, "generate N random WAN IPs")
	flag.IntVar(&scanWorkers, "w", 1024, "workers count")
	flag.IntVar(&scanWorkers, "workers", 1024, "workers count")
	flag.DurationVar(&connTimeout, "t", 700*time.Millisecond, "scan connect timeout")
	flag.DurationVar(&connTimeout, "timeout", 700*time.Millisecond, "scan connect timeout")

	flag.StringVar(&resultCallback, "cb", "", "callback to run as shell command. Use {result} as template")
	flag.DurationVar(&resultCallbackTimeout, "cbt", 30*time.Second, "callback timeout")
	flag.IntVar(&resultCallbackConcurrency, "cbmc", 30, "callback max concurrency")
	flag.BoolVar(&resultCallbackI, "cbdi", false, "disable callback infos")
	flag.BoolVar(&resultCallbackW, "cbdw", false, "disable callback warnings")
	flag.BoolVar(&resultCallbackE, "cbde", false, "disable callback errors")

	flag.StringVar(&scanPorts, "p", "", "scan ports on every rarget")
	flag.StringVar(&scanPorts, "ports", "", "scan ports on every rarget")

	flag.BoolVar(&scanRtsp, "rtsp", false, "check rtsp")
	flag.StringVar(&rtspFuzzDict, "rtspd", "./data/rtsp-paths.txt", "RTSP dictionary to fuzz")

	flag.StringVar(&cpuprofile, "profcpu", "", "profile cpu")
	flag.StringVar(&memprofile, "profmem", "", "profile mem")
}

func main() {
	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			utils.EPrintln("[E]", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			utils.EPrintln("[E]", err)
			os.Exit(1)
		}
		defer f.Close()
		defer pprof.WriteHeapProfile(f)
	}

	ipGenerator := generators.NewIPGenerator(128, randomIPsCount)
	processor := services.NewProcessor(ipGenerator, scanWorkers)

	var cbFlags utils.Flags
	if !resultCallbackE {
		cbFlags = cbFlags.Set(utils.ERR)
	}
	if !resultCallbackW {
		cbFlags = cbFlags.Set(utils.WARN)
	}
	if !resultCallbackI {
		cbFlags = cbFlags.Set(utils.INFO)
	}

	if scanRtsp {
		os.Stderr.WriteString(fmt.Sprintln("[i] Using rtsp"))
		paths, err := loaders.FileToArray(rtspFuzzDict)
		if err != nil {
			utils.EPrintln("[E]", err)
			os.Exit(1)
		}
		rtspService := services.NewRTSPService(554, paths, connTimeout)
		processor.AddService(rtspService)
	}

	if scanPorts != "" {
		ports := utils.ParseRange(scanPorts)
		processor.AddService(services.NewPortscanService(ports, connTimeout))
	}

	if len(processor.Services()) == 0 {
		processor.AddService(services.NewDummyService())
	}

	wg := new(sync.WaitGroup)
	guard := make(chan struct{}, resultCallbackConcurrency)

	for result := range processor.Process() {
		if resultCallback != "" {
			wg.Add(1)
			guard <- struct{}{}
			cmd := result.ReplaceVars(resultCallback)
			go func() {
				utils.RunCommand(cmd, wg, resultCallbackTimeout, cbFlags)
				<-guard
			}()
		} else {
			fmt.Println(&result)
		}
	}

	wg.Wait()
}
