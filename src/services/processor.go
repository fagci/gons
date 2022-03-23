package services

import (
	"go-ns/src/generators"
)

type Processor struct {
	WorkersCount int
	generator    *generators.IPGenerator
	services     []Service
	ch           chan string
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

func (p *Processor) Process() <-chan string {
	p.ch = make(chan string)
    // TODO: close channel when generator done

	for i := 0; i < p.WorkersCount; i++ {
		go p.work()
	}

	return p.ch
}

func (p *Processor) work() {
	for ip := range p.generator.GenerateWAN() {
		for _, svc := range p.services {
			for result := range svc.Check(ip.String()) {
				p.ch <- result
			}
		}
	}
}
