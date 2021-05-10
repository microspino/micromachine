package micromachine

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type assertion struct {
	have []string
	want []string
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

var fsm MicroMachine
var current string

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	fsm = *NewMicroMachine("pending")

	fsm.When("confirm", map[string]string{"pending": "confirmed"})
	fsm.When("ignore", map[string]string{"pending": "ignored"})
	fsm.When("reset", map[string]string{"confirmed": "pending", "ignored": "pending"})

	fsm.On("pending", func(event string, payload ...string) { fsm.State = "pending" })
	fsm.On("confirmed", func(event string, payload ...string) { fsm.State = "confirmed" })

	fsm.On("any", func(event string, payload ...string) {
		current = fmt.Sprintf("---> %s", fsm.State)
	})

	os.Exit(m.Run())
}

func TestIntrospection(t *testing.T) {
	tests := []assertion{
		assertion{have: fsm.Events(), want: []string{"confirm", "ignore", "reset"}},
		assertion{have: fsm.TriggerableEvents(), want: []string{"confirm", "ignore"}},
		assertion{have: fsm.States(), want: []string{"confirmed", "ignored", "pending"}},
	}

	var have, want string
	for _, test := range tests {
		have = strings.Join(test.have, ",")
		want = strings.Join(test.want, ",")
		equals(t, have, want)
	}
}

func TestExecutesCallbacksEnteringState(t *testing.T) {
	fsm.Trigger("confirm")
	want := "confirmed"
	equals(t, fsm.State, want)

	fsm.Trigger("reset")
	want = "pending"
	equals(t, fsm.State, want)
}

func TestExecutesCallbacksOnAnyTransition(t *testing.T) {
	fsm.Trigger("confirm")
	want := "---> confirmed"
	equals(t, current, want)

	fsm.Trigger("reset")
	want = "---> pending"
	equals(t, current, want)
}

func TestPassingTheEventNameToTheCallbacks(t *testing.T) {
	rcvEvent := "UNKNOWN"
	machine := NewMicroMachine("pending")
	machine.When("kill", map[string]string{"pending": "killed"})

	machine.On("killed", func(evt string, pay ...string) {
		rcvEvent = evt
	})

	machine.Trigger("kill")

	equals(t, rcvEvent, "kill")
}

func TestPassingThePayloadFromTransitionToTheCallbacks(t *testing.T) {
	rcvPayload := "UNKNOWN"
	machine := NewMicroMachine("pending")
	machine.When("implode", map[string]string{"pending": "imploded"})

	machine.On("imploded", func(evt string, pay ...string) {
		rcvPayload = pay[0]
	})

	implodePayload := "payload: tons of pressure"

	machine.Trigger("implode", implodePayload)

	equals(t, rcvPayload, implodePayload)
}
