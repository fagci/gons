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
	"strings"
	"sync"
	"time"
)

var (
	randomIPsCount                   int64
	scanWorkers, callbackConcurrency int
	connTimeout, callbackTimeout     time.Duration
	callback                         string
	callbackI, callbackE, callbackW  bool
	scanPorts                        string
	service                          string
	rtspPort                         int
	fuzzDict                         string
	cpuprofile, memprofile           string
)

func init() {
	flag.Int64Var(&randomIPsCount, "n", -1, "generate N random WAN IPs")
	flag.IntVar(&scanWorkers, "w", 1024, "workers count")
	flag.IntVar(&scanWorkers, "workers", 1024, "workers count")
	flag.DurationVar(&connTimeout, "t", 700*time.Millisecond, "scan connect timeout")
	flag.DurationVar(&connTimeout, "timeout", 700*time.Millisecond, "scan connect timeout")

	flag.StringVar(&callback, "cb", "", "callback to run as shell command. Use {result} as template")
	flag.StringVar(&callback, "callback", "", "callback to run as shell command. Use {result} as template")
	flag.DurationVar(&callbackTimeout, "cbt", 30*time.Second, "callback timeout")
	flag.IntVar(&callbackConcurrency, "cbmc", 30, "callback max concurrency")
	flag.BoolVar(&callbackI, "cbdi", false, "disable callback infos")
	flag.BoolVar(&callbackW, "cbdw", false, "disable callback warnings")
	flag.BoolVar(&callbackE, "cbde", false, "disable callback errors")

	flag.StringVar(&scanPorts, "p", "", "scan ports on every rarget")
	flag.StringVar(&scanPorts, "ports", "", "scan ports on every rarget")

	flag.StringVar(&service, "s", "", "check service (rtsp, ...)")
	flag.StringVar(&fuzzDict, "d", "./data/rtsp-paths.txt", "dictionary to fuzz")

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
	if !callbackE {
		cbFlags = cbFlags.Set(utils.ERR)
	}
	if !callbackW {
		cbFlags = cbFlags.Set(utils.WARN)
	}
	if !callbackI {
		cbFlags = cbFlags.Set(utils.INFO)
	}

	if service != "" {
		var paths []string
		var err error
		if fuzzDict != "" {
			paths, err = loaders.FileToArray(fuzzDict)
			if err != nil {
				utils.EPrintln("[E]", err)
				os.Exit(1)
			}
		}

		ports := utils.ParseRange(scanPorts)

		utils.EPrintln("[i] Using", service)
		var svc services.Service
		switch strings.ToLower(service) {
		case "rtsp":
			svc = services.NewRTSPService(ports, connTimeout, paths)
        case "portscan":
            svc = services.NewPortscanService(ports, connTimeout)
		}
		processor.AddService(svc)
	}

	if len(processor.Services()) == 0 {
		processor.AddService(services.NewDummyService())
	}

	sp := utils.Spinner{}
	sp.Start()
	defer sp.Stop()

	results := processor.Process()

	if callback != "" {
		wg := new(sync.WaitGroup)
		guard := make(chan struct{}, callbackConcurrency)
		for result := range results {
			wg.Add(1)
			guard <- struct{}{}
			cmd := result.ReplaceVars(callback)
			go func() {
				sp.Stop()
				utils.RunCommand(cmd, wg, callbackTimeout, cbFlags)
				<-guard
				if len(guard) == 0 {
					sp.Start()
				}
			}()
		}
		wg.Wait()
	} else {
		for result := range results {
			sp.Stop()
			fmt.Println(&result)
			sp.Start()
		}
	}
}
