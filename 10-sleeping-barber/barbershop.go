package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarberDoneChan  chan bool
	ClientChan      chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for client.", barber)

		for {
			// if there are no clients, the barber goes to sleep
			if len(shop.ClientChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}

				// cut hair
				shop.cutHair(client, barber)
			} else {
				//shop is closed, to send the barber home and close this goroutine
				shop.sendBarberHome(barber)
				return
			}

		}
	}()
}

func (shop *BarberShop) cutHair(client, barber string) {
	color.Green("%s is cutting %s's hair", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished %s's hair", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarberDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day.")

	close(shop.ClientChan)
	shop.Open = false

	for i := 1; i < shop.NumberOfBarbers; i++ {
		<-shop.BarberDoneChan
	}

	close(shop.BarberDoneChan)

	color.Green("-------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day and everyone has gone home")
}

func (shop *BarberShop) addClient(client string) {
	// print out message
	color.Green("*** %s arrives!", client)

	if shop.Open {
		select {
		case shop.ClientChan <- client:
			color.Blue("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Red("The shop is already closed, so %s leaves!", client)
	}
}
