package output

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type Status struct {
	info         statusInfo
	quit, exited chan bool
	m            *sync.Mutex
	w            io.Writer
}

type statusInfo struct {
	Steps    int    `json:"steps"`
	Step     int    `json:"step"`
	StepName string `json:"step-name"`
	Parts    int    `json:"parts"`
	Part     int    `json:"part"`
}

func NewStatus(totalSteps int, outWriter io.Writer) *Status {
	s := &Status{
		info:   statusInfo{
			Steps:    totalSteps,
			Step:     0,
			StepName: "Initializing",
			Parts:    0,
			Part:     0,
		},
		quit:   make(chan bool, 1),
		exited: make(chan bool, 1),
		w:      outWriter,
	}

	go func() {
		seconds := time.NewTicker(time.Second).C
		for {
			select {
			case <-seconds:
				s.print()
			case <-s.quit:
				s.print()
				s.exited <- true
				return
			}
		}
	}()
	return s
}

func (s *Status) NewStep(name string, parts int) {
	s.m.Lock()
	s.info.StepName = name
	s.info.Step++
	s.info.Parts = parts
	s.info.Part = 0
	s.m.Unlock()
}

func (s *Status) AddPart() {
	s.m.Lock()
	s.info.Part++
	s.m.Unlock()
}

func (s *Status) Stop() {
	s.m.Lock()
	s.quit <- true
	<- s.exited // Wait for exited to return something
	s.m.Unlock()
}

func (s *Status) print() {
	s.m.Lock()
	data, _ := json.Marshal(s.info)
	_, _ = s.w.Write(data)
	s.m.Unlock()
}
