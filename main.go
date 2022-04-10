package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fagci/gonr/generators"
	"github.com/fagci/gons/loaders"
	"github.com/fagci/gons/network"
	"github.com/fagci/gons/services"
	"github.com/fagci/gons/utils"
)

var (
	iface          string
	randomIPsCount int64
	cidrNetwork    string
	ipList         string
	scanWorkers    int
	connTimeout    time.Duration
	scanPorts      string
	service        string
	fuzzDict       string
)

var (
	headerReg, bodyReg string
)

var (
	callback                        string
	callbackTimeout                 time.Duration
	callbackConcurrency             int
	callbackI, callbackE, callbackW bool
)

func init() {
	flag.StringVar(&iface, "i", "", "use specific network interface")
	flag.Int64Var(&randomIPsCount, "n", -1, "generate N random WAN IPs")
	flag.StringVar(&cidrNetwork, "net", "", "Network in CIDR notation to scan in random order")
	flag.StringVar(&ipList, "list", "", "IP/networks list (CIDR) to scan in random order")
	flag.IntVar(&scanWorkers, "w", 64, "workers count")
	flag.IntVar(&scanWorkers, "workers", 64, "workers count")
	flag.DurationVar(&connTimeout, "t", 700*time.Millisecond, "scan connect timeout")
	flag.DurationVar(&connTimeout, "timeout", 700*time.Millisecond, "scan connect timeout")

	flag.StringVar(&scanPorts, "p", "", "scan ports on every rarget")
	flag.StringVar(&scanPorts, "ports", "", "scan ports on every rarget")

	flag.StringVar(&service, "s", "", "check service (rtsp, ...)")
	flag.StringVar(&fuzzDict, "d", "./assets/data/rtsp-paths.txt", "dictionary to fuzz")

	flag.StringVar(&headerReg, "rh", "", "Regexp for header")
	flag.StringVar(&bodyReg, "rb", "", "Regexp for body")

	flag.StringVar(&callback, "cb", "", "callback to run as shell command. Use {result} as template")
	flag.StringVar(&callback, "callback", "", "callback to run as shell command. Use {result} as template")
	flag.DurationVar(&callbackTimeout, "cbt", 30*time.Second, "callback timeout")
	flag.IntVar(&callbackConcurrency, "cbmc", 30, "callback max concurrency")
	flag.BoolVar(&callbackI, "cbdi", false, "disable callback infos")
	flag.BoolVar(&callbackW, "cbdw", false, "disable callback warnings")
	flag.BoolVar(&callbackE, "cbde", false, "disable callback errors")
}

func setupSercices(processor *services.Processor) {
	if service == "" {
		processor.AddService(services.NewDummyService())
	} else {
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

		var svc services.ServiceInterface

		switch strings.ToLower(service) {
		case "http":
			svc = services.NewHTTPService(ports, connTimeout, paths, headerReg, bodyReg)
		case "rtsp":
			svc = services.NewRTSPService(ports, connTimeout, paths)
		case "portscan":
			svc = services.NewPortscanService(ports, connTimeout)
		}

		if svc != nil {
			utils.EPrintln("[i] Using", service)
			utils.EPrintln("[i] Workers", scanWorkers)
			if randomIPsCount > 0 {
				utils.EPrintln("[i] Random IPs count", randomIPsCount)
			}
			processor.AddService(svc)
		}
	}
}

func process(processor *services.Processor) {
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
			go func(cmd string) {
				sp.Stop()
				utils.RunCommand(cmd, wg, callbackTimeout, cbFlags)
				<-guard
				if len(guard) == 0 {
					sp.Start()
				}
			}(cmd)
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

func main() {
	flag.Parse()

	if iface != "" {
		if err := network.SetInterface(iface); err != nil {
			utils.EPrintln("[E] iface", err)
			return
		}
		utils.EPrintln("[i] Iface", iface)
	}

	var ipSource <-chan net.IP
	if ipList != "" {
		list, err := utils.LoadInput(ipList)
		if err != nil {
			utils.EPrintln("[E] IP list", err)
			return
		}
		ipSource = generators.RandomHostsFromListGen(list)
	} else if cidrNetwork == "" {
		ipGenerator := generators.NewIPGenerator(512, randomIPsCount)
		ipSource = ipGenerator.Generate()
	} else {
		ipSource = generators.RandomHostsFromCIDRGen(cidrNetwork)
	}

	processor := services.NewProcessor(ipSource, scanWorkers)

	setupSercices(processor)
	process(processor)
}
