package svc

import "go_ns/src/gen"

type Processor struct {
	WorkersCount int
	generator    *gen.IPGenerator
	services     []Service
	ch           chan Result
}

func NewProcessor(generator *gen.IPGenerator, workersCount int) *Processor {
	return &Processor{
		WorkersCount: workersCount,
		generator:    generator,
	}
}

func (p *Processor) AddService(svc Service) {
	p.services = append(p.services, svc)
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

func (p *Processor) Process() <-chan Result {
	for i := 0; i < p.WorkersCount; i++ {
		go p.work()
	}

	return p.ch
}
