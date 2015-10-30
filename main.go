package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var wga sync.WaitGroup

type alarm struct {
	bool
	*sync.Mutex
}

func (a *alarm) Set(value bool) {
	defer a.Unlock()
	a.Lock()
	a.bool = value
}

type person struct {
	Name    string
	AllDone bool
	*sync.Mutex
}

func (p *person) setAllDone(value bool) {
	defer p.Unlock()
	p.Lock()
	p.AllDone = value
}

type people []person

var alrm = alarm{false, &sync.Mutex{}}
var hsld = []person{}

func main() {
	fmt.Println("Let's go for walk!")
	Bob := person{"Bob", false, &sync.Mutex{}}
	Alice := person{"Alice", false, &sync.Mutex{}}
	hsld = []person{Bob, Alice}

	for k := range hsld {
		wg.Add(1)
		go hsld[k].dotask("getting ready", 1, 3, false)
	}

	wg.Wait()

	wga.Add(1)
	go setAlarm(5)

	for k := range hsld {
		wg.Add(1)
		go hsld[k].dotask("putting on shoes", 1, 3, true)
	}
	wg.Wait()
	alrm.Lock()
	if alrm.bool {
		fmt.Println("Crap!  The alarm is already set.")
	} else {
		fmt.Println("Exiting and locking door.")
	}
	alrm.Unlock()
	wga.Wait()
}

func setAlarm(delay int) {
	defer wga.Done()
	fmt.Println("Arming alarm.")
	fmt.Println("Alarm is counting down.")
	time.Sleep(time.Duration(delay) * time.Second)
	alrm.Set(true)
	fmt.Println("Alarm armed.")
	for _, v := range hsld {
		v.Lock()
		if v.AllDone == false {
			fmt.Printf("Alarm set before %v was ready.\n", v.Name)
		}
		v.Unlock()
	}
}

func (p *person) dotask(task string, min int, max int, setdone bool) {
	defer wg.Done()
	s := random(min, max)
	fmt.Println(p.Name, "started", task)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Println(p.Name, "spent", s, "seconds", task)
	if setdone {
		p.setAllDone(true)
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	time.Sleep(10 * time.Microsecond)
	return rand.Intn(max-min) + min
}
