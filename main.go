package main

import (
	"fmt"
	"time"
)

/*

value -> observable -> observable -> observer

channel + goroutine
*/

type Observer interface {
	OnNext(v int)
	OnError(e error)
	OnCompleted()
}

type Observable interface {
	Subscribe(Observer)
	Unsubscribe(Observer)
	Connect()
}

func fromChan(ch chan int) Observable {
	return &Stream{ch: ch}
}

type Stream struct {
	ch        chan int
	observers []Observer
}

func (s *Stream) Connect() {
	go func() {
		defer func() {
			e := recover()

			if e == nil {
				return
			}

			switch err := e.(type) {
			case error:
				for _, ob := range s.observers {
					ob.OnError(err)
				}
			default:
				for _, ob := range s.observers {
					ob.OnError(fmt.Errorf("%v", err))
				}
			}

		}()

		for v := range s.ch {
			for _, ob := range s.observers {
				ob.OnNext(v)
			}
		}
		for _, ob := range s.observers {
			ob.OnCompleted()
		}
		return
	}()
}

func (s *Stream) Subscribe(ob Observer) {
	s.observers = append(s.observers, ob)
}

func (s *Stream) Unsubscribe(observer Observer) {
	for i, ob := range s.observers {
		if ob == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return
		}
	}
}

type IntObserver struct{}

func (ob IntObserver) OnNext(v int) {
	fmt.Println("Value", v)
}
func (ob IntObserver) OnError(e error) {
	fmt.Println("Error", e)
}
func (ob IntObserver) OnCompleted() {
	fmt.Println("Completed")
}

func emitItems(ch chan int) {
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
}

func main() {
	ch := make(chan int)

	go emitItems(ch)

	obs := fromChan(ch)
	ob := IntObserver{}
	obs.Subscribe(ob)
	obs.Connect()

	time.Sleep(1000 * time.Second)
}
