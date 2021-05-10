package main

import (
	"fmt"
	"micromachine"
	"strings"
)

func main() {
	fsm := micromachine.NewMicroMachine("pending")

	fsm.When("confirm", map[string]string{"pending": "confirmed"})
	fsm.When("ignore", map[string]string{"pending": "ignored"})
	fsm.When("reset", map[string]string{"confirmed": "pending", "ignored": "pending"})

	fmt.Println("A new machine is Ready!")

	// Let's create a callback to capitalize the machine state
	capitalize := func(evt string, pay ...string) { fmt.Println(strings.Title(fsm.State)) }
	fsm.On("any", capitalize)

	fmt.Println("Should print Confirmed, Pending and Ignored...")

	fsm.Trigger("confirm")
	fsm.Trigger("ignore")
	fsm.Trigger("reset")
	fsm.Trigger("ignore")

	fmt.Println("Should print all possible states: confirmed, ignored and pending")
	fmt.Println("States -> ", fsm.States())

	fmt.Println("Should print all events: confirm, ignore, reset")
	fmt.Println("Events -> ", fsm.Events())
}
