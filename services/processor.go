package services

import (
	"github.com/fagci/gons/generators"
	"net"
	"sync"
)

type Processor struct {
	WorkersCount int
	generator    *generators.IPGenerator
	services     []Service
	ch           chan HostResult
}

func NewProcessor(generator *generators.IPGenerator, workersCount int) *Processor {
	return &Processor{
		WorkersCount: workersCount,
		generator:    generator,
	}
}

func (p *Processor) AddService(svc Service) {
	p.services = append(p.services, svc)
}

func (p *Processor) Services() []Service {
	return p.services
}

func (p *Processor) Process() <-chan HostResult {
	var wg sync.WaitGroup
	p.ch = make(chan HostResult)

	ipGeneratorChannel := p.generator.Generate()
	for i := 0; i < p.WorkersCount; i++ {
		wg.Add(1)
		go p.work(ipGeneratorChannel, &wg)
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
