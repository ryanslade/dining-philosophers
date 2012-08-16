package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Philosopher struct {
	id    int
	right chan bool
	left  chan bool
}

func NewPhilosopher(id int) *Philosopher {
	return &Philosopher{
		id:    id,
		right: make(chan bool),
		left:  make(chan bool),
	}
}

func (p *Philosopher) Println(s string) {
	fmt.Printf("[%v] %v\n", p.id, s)
}

func sleepDuration() time.Duration {
	return time.Duration(rand.Intn(1000)) * time.Millisecond
}

func (p *Philosopher) tryEat(haveFork, waitFork chan bool) {
	p.Println("Have fork...")
	select {
	case <-waitFork:
		p.Println("Eating...")
		time.Sleep(sleepDuration())
		waitFork <- true
		haveFork <- true
	case <-time.After(sleepDuration()):
		p.Println("Giving up fork")
		haveFork <- true
	}
}

func (p *Philosopher) Dine() {
	for {
		p.Println("Thinking...")
		time.Sleep(sleepDuration())

		select {
		case <-p.right:
            p.tryEat(p.right, p.left)
		case <-p.left:
            p.tryEat(p.left, p.right)
		}

	}
}

type Place struct {
	id      int
	hasFork bool
	left    chan bool
	right   chan bool
}

func (p *Place) Wait() {
	for {
		if p.hasFork {
			select {
			case p.left <- true:
				p.hasFork = false
			case p.right <- true:
				p.hasFork = false
			}
		} else {
			select {
			case <-p.left:
				p.hasFork = true
			case <-p.right:
				p.hasFork = true
			}
		}
	}
}

const numPhilosophers = 5

func main() {
	philosophers := make([]*Philosopher, numPhilosophers)
	places := make([]*Place, numPhilosophers)

	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = NewPhilosopher(i)
	}

	for i := 0; i < numPhilosophers-1; i++ {
		places[i] = &Place{id: i, hasFork: true, left: philosophers[i].right, right: philosophers[i+1].left}
	}
	places[numPhilosophers-1] = &Place{id: numPhilosophers - 1, hasFork: true, left: philosophers[numPhilosophers-1].right, right: philosophers[0].left}

	for i := 0; i < numPhilosophers; i++ {
		go philosophers[i].Dine()
		go places[i].Wait()
	}

	select {}
}
