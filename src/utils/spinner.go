package utils

import (
	"fmt"
	"time"
)

type Spinner struct {
	chars []string
	i     int
	stop  chan struct{}
    running bool
}

const FRAME_DUR = 90 * time.Millisecond

func (s *Spinner) Start() {
    s.stop = make(chan struct{})
	s.chars = []string{"|", "/", "-", "\\"}
	s.i = 0
    s.running = true
	go s.spin()
}

func (s *Spinner) Stop() {
    if !s.running {
        return
    }
    s.running = false
    s.clear()
    defer close(s.stop)
	s.stop <- struct{}{}
}

func (s *Spinner) next() {
	fmt.Printf("\r%s", s.chars[s.i])
	s.i++
	if s.i >= len(s.chars) {
		s.i = 0
	}
}
func (s *Spinner) clear() {
    fmt.Printf("\r")
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