package micromachine

import (
	"fmt"
	"sort"
)

type MicroMachine struct {
	State          string
	TransitionsFor map[string]map[string]string
	Callbacks      map[string][]Callback
}

type Transitions map[string]string

type Callback func(event string, payload ...string)

func NewMicroMachine(initialState string) *MicroMachine {
	m := MicroMachine{State: initialState}
	m.TransitionsFor = map[string]map[string]string{}
	m.Callbacks = map[string][]Callback{}
	return &m
}

func (m *MicroMachine) On(key string, c Callback) {
	m.Callbacks[key] = append(m.Callbacks[key], c)
}

func (m *MicroMachine) When(event string, transitions map[string]string) {
	m.TransitionsFor[event] = transitions
}

func (m *MicroMachine) canTrigger(event string) bool {
	if _, ok := m.TransitionsFor[event]; ok {
		_, ok := m.TransitionsFor[event][m.State]
		return ok
	}
	return false
}

func (m *MicroMachine) Trigger(event string, payload ...string) bool {
	safePayload := ""
	if len(payload) > 0 {
		safePayload = payload[0]
	}
	if m.canTrigger(event) {
		return m.change(event, safePayload)
	} else {
		return false
	}
}

func (m *MicroMachine) TriggerOrRaise(event string, payload ...string) error {
	safePayload := ""
	if len(payload) > 0 {
		safePayload = payload[0]
	}
	if m.Trigger(event, safePayload) {
		return nil
	} else {
		return fmt.Errorf("Event '%s' not valid from state '%s'", event, m.State)
	}
}

func (m *MicroMachine) Events() []string {
	keys := make([]string, 0, len(m.TransitionsFor))
	for k := range m.TransitionsFor {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *MicroMachine) TriggerableEvents() (triggerables []string) {
	for _, e := range m.Events() {
		if m.canTrigger(e) {
			triggerables = append(triggerables, e)
		}
	}
	sort.Strings(triggerables)
	return
}

func (m *MicroMachine) States() []string {
	in := []string{}
	for _, transitions := range m.TransitionsFor {
		for a, b := range transitions {
			in = append(in, a, b)
		}
	}
	sort.Strings(in)
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		in[j] = in[i]
	}
	states := in[:j+1]
	return states
}

func (m *MicroMachine) change(event, payload string) bool {
	m.State = m.TransitionsFor[event][m.State]
	var callbacks []Callback
	callbacks = append(callbacks, m.Callbacks["any"]...)
	callbacks = append(callbacks, m.Callbacks[m.State]...)
	for _, c := range callbacks {
		c(event, payload)
	}
	return true
}
