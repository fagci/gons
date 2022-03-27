package services

import (
	"gons/src/generators"
	"gons/src/models"
	"net"
	"sync"
)

type Processor struct {
	WorkersCount int
	generator    *generators.IPGenerator
	services     []Service
	ch           chan models.HostResult
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

func (p *Processor) Process() <-chan models.HostResult {
	var wg sync.WaitGroup
	p.ch = make(chan models.HostResult)

	ipGeneratorChannel := p.generator.GenerateWAN()
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
	for ip := range ipGeneratorChannel {
		for _, svc := range p.services {
			for result := range svc.Check(ip) {
				p.ch <- result
			}
		}
	}
}
