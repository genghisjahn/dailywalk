package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var wga sync.WaitGroup

type person struct {
	Name    string
	AllDone bool
}

var people []person
var mtx = &sync.Mutex{}
var alarmset bool

func main() {
	fmt.Println("Let's go for walk!")
	Bob := person{"Bob", false}
	Alice := person{"Alice", false}
	people = []person{Bob, Alice}

	for _, v := range people {
		wg.Add(1)
		go v.dotask("getting ready", 10, 20, false)
	}
	wg.Wait()

	wga.Add(1)
	go setAlarm(60)

	for _, v := range people {
		wg.Add(1)
		go v.dotask("putting on shoes", 10, 20, true)
	}
	wg.Wait()
	if alarmset {
		fmt.Println("Crap!  The alarm is already set.")
	} else {
		fmt.Println("Exiting and locking door.")
	}
	wga.Wait()
}

func setAlarm(delay int) {
	fmt.Println("Arming alarm.")
	fmt.Println("Alarm is counting down.")
	time.Sleep(time.Duration(delay) * time.Second)
	alarmset = true
	fmt.Println("Alarm armed.")
	mtx.Lock()
	for _, v := range people {
		if v.AllDone == false {
			fmt.Printf("Alarm set before %v was ready.\n", v.Name)
		}
	}
	defer func() {
		mtx.Unlock()
		wga.Done()
	}()
}

func (p *person) dotask(task string, min int, max int, setdone bool) {
	defer wg.Done()
	s := random(min, max)
	fmt.Println(p.Name, "started", task)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Println(p.Name, "spent", s, "seconds", task)
	if setdone {
		mtx.Lock()
		p.AllDone = true
		mtx.Unlock()
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	time.Sleep(10 * time.Microsecond)
	return rand.Intn(max-min) + min
}
