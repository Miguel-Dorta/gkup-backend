package output

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type Status struct {
	info         statusInfo
	queue        *linkedList
	quit, exited chan bool
	w            io.Writer
	printMutex   *sync.Mutex
	m            *sync.Mutex
}

type statusInfo struct {
	Steps           int    `json:"steps"`
	Step            int    `json:"step"`
	StepName        string `json:"step-name"`
	Parts           int    `json:"parts"` // Can be 0
	Part            int    `json:"part"`
}

func NewStatus(totalSteps, outputTimeInMS int, outWriter io.Writer) *Status {
	if totalSteps < 1 {
		panic("totalSteps must be >= 1 in status")
	}
	if outputTimeInMS < 1 {
		panic("invalid outputTime for status")
	}

	s := &Status{
		info: statusInfo{
			Steps:    totalSteps,
			StepName: "Initializing",
		},
		queue:  new(linkedList),
		quit:   make(chan bool, 1),
		exited: make(chan bool, 1),
		w:      outWriter,
	}

	go func() {
		outputTicker := time.NewTicker(time.Duration(outputTimeInMS * int(time.Millisecond))).C
		for {
			select {
			case <-outputTicker:
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
	s.queue.Push(s.info)
	s.m.Unlock()
}

func (s *Status) AddPart() {
	s.m.Lock()
	s.info.Part++
	s.queue.Push(s.info)
	s.m.Unlock()
}

func (s *Status) print() {
	s.m.Lock()
	n := s.queue.PopAndReset()
	s.m.Unlock()

	s.printMutex.Lock() // This is not perfect, but it's improbable that something bad happens. Let's hope!
	for ; n != nil; n = n.next {
		data, _ := json.Marshal(n.value)
		_, _ = s.w.Write(data)
	}
	s.printMutex.Unlock()
}

func (s *Status) Stop() {
	s.m.Lock()
	s.quit <- true
	<-s.exited // Wait for exited channel to return something
	s.m.Unlock()
}
