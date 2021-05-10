package main

import (
	"fmt"
	"micromachine"
)

func main() {
	fsm := micromachine.NewMicroMachine("pending")

	fsm.When("confirm", map[string]string{"pending": "confirmed"})
	fsm.When("ignore", map[string]string{"pending": "ignored"})
	fsm.When("reset", map[string]string{"confirmed": "pending", "ignored": "pending"})

	fmt.Println("A new machine is Ready!")

	fmt.Println("Should print: confirmed, pending and ignored...")

	// pending -> confirmed
	if fsm.Trigger("confirm") {
		fmt.Println("-> ", fsm.State)
	}

	// confirmed -> ignored transition cannot be triggered
	// so the follwing will not do/print anything
	// the fsm cannot be "ignored" from confirmed
	if fsm.Trigger("ignore") {
		fmt.Println("-> ", fsm.State)
	}

	// pending <- confirmed
	if fsm.Trigger("reset") {
		fmt.Println("-> ", fsm.State)
	}

	// pending -> ignored
	if fsm.Trigger("ignore") {
		fmt.Println("-> ", fsm.State)
	}
}
