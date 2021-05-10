package main

import (
	"fmt"
	"micromachine"
)

type transitions map[string]string

func buildOrder() *micromachine.MicroMachine {
	o := micromachine.NewMicroMachine("start")

	o.When("create", transitions{"start": "payment"})
	o.When("pay", transitions{"payment": "shipment"})
	o.When("ship", transitions{"shipment": "shipped"})
	o.When("cancel", transitions{"shipment": "awaiting-refund", "payment": "canceled"})
	o.When("refund", transitions{"awaiting-refund": "canceled"})

	displayState := func(e string, p ...string) { fmt.Println(o.State) }

	for _, state := range o.States() {
		o.On(state, displayState)
	}

	o.On("any", func(e string, payload ...string) {
		fmt.Printf("⚡️[%s]\n", e)
	})

	return o
}

func main() {
	o1 := buildOrder()
	fmt.Println("States 🏠 ", o1.States())
	fmt.Println("Events ⚡️ ", o1.Events())

	// simple buy and ship
	fmt.Println("\nLet's order something...")
	fmt.Println(o1.State)
	o1.Trigger("create")
	o1.Trigger("pay")
	o1.Trigger("ship")

	// refund
	o2 := buildOrder()
	fmt.Println("\nSomeone asked for refund...")
	fmt.Println(o2.State)
	o2.Trigger("create")
	o2.Trigger("pay")
	o2.Trigger("cancel")
	o2.Trigger("refund")
}
