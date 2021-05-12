package main

import (
	"fmt"
	MM "micromachine"
)

func buildOrder() *MM.MicroMachine {
	o := micromachine.NewMicroMachine("start")

	o.When("create", MM.Transitions{"start": "payment"})
	o.When("pay", MM.Transitions{"payment": "shipment"})
	o.When("ship", MM.Transitions{"shipment": "shipped"})
	o.When("cancel", MM.Transitions{"shipment": "wait-refund", "payment": "canceled"})
	o.When("refund", MM.Transitions{"wait-refund": "canceled"})

	o.On("any", func(evt string, payload ...string) {
		if payload[0] != "" {
			evt += fmt.Sprintf(" 🎒%s", payload[0])
		}
		fmt.Printf("↓ %s\n·%s\n", evt, o.State)
	})

	return o
}

func main() {
	o1 := buildOrder()
	fmt.Println("States ·", o1.States())
	fmt.Println("Events ↓", o1.Events())

	// simple buy and ship
	fmt.Println("\nLet's order something...")
	fmt.Println(o1.State)
	o1.Trigger("create")
	o1.Trigger("pay", ">>800€")
	o1.Trigger("ship")

	// refund
	o2 := buildOrder()
	fmt.Println("\nSomeone asked for refund...")
	fmt.Println(o2.State)
	o2.Trigger("create")
	o2.Trigger("pay", ">>300€")
	o2.Trigger("cancel")
	o2.Trigger("refund", "<<300€")
}
