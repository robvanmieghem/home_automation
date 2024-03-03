// Example listening for ptm216b ble events on Darwin
package main

import (
	"context"
	"log"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
	"github.com/robvanmieghem/home_automation/enocean/ptmb"
)

func AdvertisementHandler(advertisement ble.Advertisement) {
	event, err := ptmb.NewEvent(advertisement)
	if err != nil {
		log.Fatal(err)
	}
	if event == nil { // No EnOcean ptm216b event
		return
	}

	println("EnOcean event received")
	println("Address:", event.Address)
	println("Sequence:", event.Sequence)
	switch true {
	case event.State.IsButtonA0():
		println("Button A0")

	case event.State.IsButtonA1():
		println("Button A1")

	case event.State.IsButtonB0():
		println("Button B0")

	case event.State.IsButtonB1():
		println("Button B1")
	}

}

func main() {
	d, err := darwin.NewDevice()
	if err != nil {
		log.Fatal(err)
	}
	d.Scan(context.Background(), false, AdvertisementHandler)
}
