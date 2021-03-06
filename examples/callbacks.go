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

	displayState := func(e string, p ...string) { fmt.Println(fsm.State) }

	fsm.On("confirmed", displayState)
	fsm.On("ignored", displayState)
	fsm.On("pending", displayState)

	// Will print:
	// confirmed Transitioned... reset Transitioned... ignored Transitioned...
	fsm.On("any", func(e string, payload ...string) {
		fmt.Printf("⚡️[%s] %s\n", e, strings.Join(payload, ","))
	})

	fmt.Println(fsm.State)

	fsm.Trigger("confirm")
	fsm.Trigger("ignore")
	fsm.Trigger("reset", "payload1", "payload2", "payload3") // will use only "payload1"
	fsm.Trigger("ignore")
}
