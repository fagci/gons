package services

import (
	"go-ns/src/generators"
	"go-ns/src/models"
	"net"
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
	p.ch = make(chan models.HostResult)
	// TODO: close channel when generator done
	ch := p.generator.GenerateWAN()
	for i := 0; i < p.WorkersCount; i++ {
		go p.work(ch)
	}

	return p.ch
}

func (p *Processor) work(ipGeneratorChannel <-chan net.IP) {
	for ip := range ipGeneratorChannel {
		for _, svc := range p.services {
			for result := range svc.Check(ip) {
				p.ch <- result
			}
		}
	}
}
