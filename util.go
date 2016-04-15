package main

import (
	"errors"
	"time"
)

type result struct {
	val string
	err error
}

func loop(f func() result, periodS int, resultCh chan<- result) {
	updateCh := make(chan result)
	go func() {
		for {
			updateCh <- f()
			time.Sleep(time.Duration(periodS) * time.Second)
		}
	}()

	result := result{err: errors.New("No value yet")}
	for {
		select {
		case resultCh <- result:
		case result = <-updateCh:
		}
	}
}
