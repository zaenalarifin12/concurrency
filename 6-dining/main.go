package main

import (
	"fmt"
	"sync"
	"time"
)

// The dining philosophers problem is well known in computer science circles
// Five philosophers, numbered from 0 trough 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table
// Their only difficulty - besides those of philosophy - is the dish
// served is very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating , since there are five philosopher and five forks.

// This is a simple implementation of Dijkstra's solution to the "Dining Philosopher" dilemma

// Philosopher is a struct which stores information about a philosopher
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// Philosopher is list of all philosopher
var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Scorates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

// define some variable
var hunger = 3
var eatTime = 0 * time.Second
var thinkTime = 0 * time.Second
var sleepTime = 0 * time.Second

var orderMutex sync.Mutex // a mutext for the slice orderFinished; part of challenge
var orderFinish []string  // the order in which philosopher finish in dining and leave; part of the challenge

func main() {
	// print out a welcome message
	fmt.Println("Dining Philosophers problem")
	fmt.Println("----------------------------")
	fmt.Println("The table is empty")

	time.Sleep(sleepTime)
	// start the meal
	dine()

	// print out finish message
	time.Sleep(sleepTime)
	fmt.Println("The table is empty")
}

func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map if all 5 forks
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		// fire off a goroutine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()
}

// diningProblem is the function fired of ass goroutine for each of our philosophers. it takes one
// philosopher, our WaitGroup to determine when everyone is done, a map containing the mutexes for every
// fork on the table, and a WaitGroup used to pause execution of every instance of this goroutine
// until everyone is seated at the table

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done()

	seated.Wait()
	// eat three times
	for i := hunger; i > 0; i-- {

		if philosopher.leftFork > philosopher.rightFork {
			//right first || lower number
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)

		}
		//// get a lock on both forks
		//forks[philosopher.leftFork].Lock()
		//fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
		//forks[philosopher.rightFork].Lock()
		//fmt.Printf("\t%s takes the right fork.\n", philosopher.name)

		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks.\n", philosopher.name)

	}

	fmt.Println(philosopher.name, "is satisfied")
	fmt.Println(philosopher.name, "left the table")
	
	orderMutex.Lock()
	orderFinish = append(orderFinish, philosopher.name)
	orderMutex.Unlock()
}
