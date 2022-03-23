package services

import (
	"go-ns/src/generators"
	"go-ns/src/models"
)

type Processor struct {
	WorkersCount int
	generator    *generators.IPGenerator
	services     []Service
	ch           chan models.Result
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

func (p *Processor) Process() <-chan models.Result {
	p.ch = make(chan models.Result)
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
