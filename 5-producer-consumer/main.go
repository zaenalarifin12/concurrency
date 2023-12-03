package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var BlackString = 30
var RedString = 31
var GreenString = 32
var YellowString = 33
var BlueString = 34
var MagentaString = 35
var CyanString = 36
var WhiteString = 37
var HiBlackString = 90
var HiRedString = 91
var HiGreenString = 92
var HiYellowString = 93
var HiBlueString = 94
var HiMagentaString = 95
var HiCyanString = 96
var HiWhiteString = 97

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzaria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0
	// run forever or until we receive a quit notification
	// try to make pizzas
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// we tried to make a pizza ( we sent something to the data channel)
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				// close channel
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func printColor(number int, text string, additionalParams ...interface{}) {
	formatString := "\x1b[%dm%s\x1b[0m"

	if len(additionalParams) > 0 {
		formatString += " " // Add a space before additional parameters
	}

	formatString += strings.Repeat("%v ", len(additionalParams))

	// Corrected the order of parameters in the allParams slice
	allParams := append([]interface{}{number, text}, additionalParams...)
	fmt.Printf(formatString, allParams...)
	fmt.Println()
}

func main() {
	// seed the random number generator
	//seed := time.Now().UnixNano()
	//randSource := rand.NewSource(seed)
	//rand.New(randSource).Intn(100)

	//print out a message
	printColor(CyanString, "The pizzaria is open for bussiness")
	printColor(CyanString, "----------------------------------")

	//create producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzaria(pizzaJob)

	// create and run customer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				printColor(GreenString, i.message)
				printColor(GreenString, fmt.Sprintf("Order %d is out for delivery!", i.pizzaNumber))
			} else {
				printColor(RedString, i.message)
				printColor(RedString, "The customer is really mad!")
			}
		} else {
			printColor(CyanString, "Done making pizza....")
			err := pizzaJob.close()
			if err != nil {
				printColor(RedString, fmt.Sprintf("*** Error Closing channel! %d", err))
			}
		}
	}

	// print out the ending message
	printColor(CyanString, "-----------------")
	printColor(CyanString, "Done for the day")

	printColor(CyanString, fmt.Sprintf("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total))

	switch {
	case pizzasFailed > 9:
		printColor(RedString, "it was awful day!")
	case pizzasFailed >= 6:
		printColor(RedString, "It was not very good day!")
	case pizzasFailed >= 4:
		printColor(YellowString, "It was an okay day!")
	case pizzasFailed >= 2:
		printColor(YellowString, "It was a pretty good day!")
	default:
		printColor(GreenString, "It was a great day!")
	}
}
