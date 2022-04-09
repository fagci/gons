package services

import (
	"net"
	"sync"
)

type Processor struct {
	WorkersCount int
	ipSource     <-chan net.IP
	services     []ServiceInterface
	ch           chan HostResult
}

func NewProcessor(ipSource <-chan net.IP, workersCount int) *Processor {
	return &Processor{
		WorkersCount: workersCount,
		ipSource:     ipSource,
	}
}

func (p *Processor) AddService(svc ServiceInterface) {
	p.services = append(p.services, svc)
}

func (p *Processor) Services() []ServiceInterface {
	return p.services
}

func (p *Processor) Process() <-chan HostResult {
	var wg sync.WaitGroup
	p.ch = make(chan HostResult)

	for i := 0; i < p.WorkersCount; i++ {
		wg.Add(1)
		go p.work(p.ipSource, &wg)
	}

	go func() {
		defer close(p.ch)
		wg.Wait()
	}()

	return p.ch
}

func (p *Processor) work(ipGeneratorChannel <-chan net.IP, wg *sync.WaitGroup) {
	defer wg.Done()
	var swg sync.WaitGroup
	for ip := range ipGeneratorChannel {
		for _, svc := range p.services {
			swg.Add(1)
			go svc.Check(ip, p.ch, &swg)
		}
	}
	swg.Wait()
}
