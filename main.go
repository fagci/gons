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
	utils.EPrintln("=========================")
	utils.EPrintln("NetStalking tool by fagci")
	utils.EPrintln("-------------------------")

	flag.StringVar(&iface, "i", "", "use specific network interface")
	flag.Int64Var(&randomIPsCount, "n", -1, "generate N random WAN IPs")
	flag.StringVar(&cidrNetwork, "net", "", "Network in CIDR notation to scan in random order")
	flag.StringVar(&ipList, "list", "", "IP/networks list (CIDR) to scan in random order")
	flag.IntVar(&scanWorkers, "w", 4096, "workers count")
	flag.IntVar(&scanWorkers, "workers", 4096, "workers count")
	flag.DurationVar(&connTimeout, "t", 700*time.Millisecond, "scan connect timeout")
	flag.DurationVar(&connTimeout, "timeout", 700*time.Millisecond, "scan connect timeout")

	flag.StringVar(&scanPorts, "p", "", "scan ports on every rarget")
	flag.StringVar(&scanPorts, "ports", "", "scan ports on every rarget")

	flag.StringVar(&service, "s", "", "check service (rtsp, ...)")
	flag.StringVar(&fuzzDict, "d", "", "dictionary to fuzz")

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
		var (
			err   error
			paths []string
			svc   services.ServiceInterface
		)

		service := strings.ToLower(service)

		switch service {
		case "portscan":
			fuzzDict = ""
			break
		case "rtsp":
			if fuzzDict == "" {
				fuzzDict = "./assets/data/rtsp-paths.txt"
			}
			break
		case "http":
			if fuzzDict == "" {
				fuzzDict = "./assets/data/http-cam-paths.txt"
				headerReg = "(multipart/x-mixed-replace|image/jpeg)"
			}
			break
		}

		if fuzzDict != "" {
			paths, err = loaders.FileToArray(fuzzDict)
			if err != nil {
				utils.EPrintln("[E]", err)
				os.Exit(1)
			}
		}

		ports := utils.ParseRange(scanPorts)

		switch service {
		case "http":
			svc = services.NewHTTPService(ports, connTimeout, paths, headerReg, bodyReg)
		case "rtsp":
			svc = services.NewRTSPService(ports, connTimeout, paths)
		case "portscan":
			svc = services.NewPortscanService(ports, connTimeout)
		}

		if svc != nil {
			utils.EPrintln("service:     ", service)
			if len(ports) != 0 {
				utils.EPrintln("ports:       ", scanPorts)
			}
			utils.EPrintln("workers:     ", scanWorkers)
			utils.EPrintln("conn timeout:", connTimeout)
			processor.AddService(svc)
		}
		if randomIPsCount > 0 {
			utils.EPrintln("max hosts:   ", randomIPsCount)
		}
		if svc != nil && fuzzDict != "" {
			utils.EPrintln("dict:", fuzzDict)
		}
	}
}

func process(processor *services.Processor) {
	var cbFlags utils.Flags

	sp := utils.Spinner{}

	results := processor.Process()
	if callback == "" {
		utils.EPrintln("=========================")
		sp.Start()
		defer sp.Stop()
		for result := range results {
			sp.Stop()
			fmt.Println(&result)
			sp.Start()
		}
		return
	}

	utils.EPrint("[i] callback set")
	if !callbackE {
		utils.EPrint(" [err]")
		cbFlags = cbFlags.Set(utils.ERR)
	}
	if !callbackW {
		utils.EPrint(" [warn]")
		cbFlags = cbFlags.Set(utils.WARN)
	}
	if !callbackI {
		utils.EPrint(" [info]")
		cbFlags = cbFlags.Set(utils.INFO)
	}
	utils.EPrintln()
	utils.EPrintln("    cb concurrency:", callbackConcurrency)
	utils.EPrintln("    cb timeout:", callbackTimeout)

	utils.EPrintln("=========================")
	sp.Start()
	defer sp.Stop()

	wg := new(sync.WaitGroup)
	guard := make(chan struct{}, callbackConcurrency)
	for result := range results {
		wg.Add(1)
		guard <- struct{}{}
		sp.Stop()
		go func(cmd string) {
			utils.RunCommand(cmd, wg, callbackTimeout, cbFlags)
			if len(guard) == 1 {
				sp.Start()
			}
			<-guard
		}(result.ReplaceVars(callback))
	}
	wg.Wait()
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
