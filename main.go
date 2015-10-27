package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

type person struct {
	Name    string
	AllDone bool
}

var people []person

func main() {
	Bob := person{"Bob", false}
	Alice := person{"Alice", false}
	people = []person{Bob, Alice}

	wg.Add(2)
	go Bob.dotask("getting ready", 10, 20, false)
	go Alice.dotask("getting ready", 10, 20, false)
	wg.Wait()

	wg.Add(1)
	go setAlarm(60)

	wg.Add(2)
	go Bob.dotask("putting on shoes", 10, 20, true)
	go Alice.dotask("putting on shoes", 10, 20, true)
	fmt.Println("Exiting and locking door.")
	wg.Wait()
}

func setAlarm(delay int) {
	fmt.Println("Arming alarm.")
	fmt.Println("Alarm is counting down.")
	time.Sleep(time.Duration(delay) * time.Second)
	fmt.Println("Alarm armed.")
	for _, v := range people {
		if !v.AllDone {
			panic(fmt.Sprintf("Alarm set before %v was ready.", v.Name))
		}
	}
	defer wg.Done()
}

func (p person) dotask(task string, min int, max int, setdone bool) {
	defer wg.Done()
	s := random(min, max)
	fmt.Println(p.Name, "started", task)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Println(p.Name, "spent", s, "seconds", task)
	if setdone {
		p.AllDone = true
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	time.Sleep(10 * time.Microsecond)
	return rand.Intn(max-min) + min
}
