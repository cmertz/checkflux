package main

import (
	"fmt"
	"math/rand"
	"time"
)

type state int

const (
	stateOk state = iota
	stateNotOk
	stateUnknown
)

var stateLabels = map[state]string{
	stateOk:      "OK",
	stateNotOk:   "CRITICAL",
	stateUnknown: "UNKNOWN",
}

type result struct {
	ID    string
	state state
}

type execution struct {
	ID      string
	Check   func() (state, error)
	results chan<- result
	Timer   <-chan time.Time
}

func (e execution) Perform() {
	for {
		select {
		case <-e.Timer:
			r := result{}

			s, _ := e.Check()

			r.state = s
			r.ID = e.ID

			e.results <- r
		}
	}
}

func randomInterval() time.Duration {
	var r int
	for ; r < 2; r = rand.Intn(10) {
	}

	res, _ := time.ParseDuration(fmt.Sprintf("%ds", r))

	return res
}

func randomPercentage() int {
	return rand.Intn(11) + 89
}

func probability(percent int) bool {
	if percent > 100 {
		panic("")
	}

	r := rand.Intn(99)
	if r > percent-1 {
		return false
	}

	return true
}

func newRandom(prob int) func() (state, error) {
	return func() (state, error) {
		if probability(prob) {
			return stateOk, nil
		}
		if probability(prob) {
			return stateNotOk, nil
		}
		return stateUnknown, nil
	}
}

func newRandomCheck(id int, results chan result) *execution {
	return &execution{
		ID:      fmt.Sprintf("%d", id),
		Check:   newRandom(randomPercentage()),
		results: results,
		Timer:   time.Tick(randomInterval()),
	}
}
