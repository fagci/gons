package utils

import (
	"time"
)

type Spinner struct {
	chars   []string
	i       int
	stop    chan struct{}
	running bool
}

const FRAME_DUR = 90 * time.Millisecond

func (s *Spinner) Start() {
	if s.running {
		return
	}
	s.stop = make(chan struct{})
	s.chars = []string{"|", "/", "-", "\\"}
	s.i = 0
	s.running = true
	go s.spin()
}

func (s *Spinner) Stop() {
	if s.running {
		s.running = false
		s.clear()
		defer close(s.stop)
		s.stop <- struct{}{}
	}
}

func (s *Spinner) next() {
	EPrintf("\r%s", s.chars[s.i])
	s.i++
	if s.i >= len(s.chars) {
		s.i = 0
	}
}
func (s *Spinner) clear() {
	EPrintf("\r")
}

func (s *Spinner) spin() {
	t := time.NewTicker(FRAME_DUR)
	for {
		select {
		case <-t.C:
			s.next()
		case <-s.stop:
			return
		}
	}
}
