// Copyright 2016 the Go-Circuits Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"
)

var calls []string

func BasicEventHandler(e Event) {
	calls = append(calls, e.Target())
}

func FailingEventHandler(e Event) {
	calls = append(calls, e.Target())
	panic("Gophers in the engine room!")
}

func Test_SimpleEvent(t *testing.T) {
	calls = make([]string, 0)
	c := NewComponent()
	c.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	c.Fire(&BaseEvent{target: "foo"})
	c.Fire(&BaseEvent{target: "exit"})
	c.Run(1)
	if len(calls) != 1 {
		t.Errorf("Expected one call to the EventHandler. Got %d.", len(calls))
	}
}

func Test_EventsFIFO(t *testing.T) {
	calls = make([]string, 0)
	c := NewComponent()
	c.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	c.RegisterEventHandler(NewEventHandler("bar", BasicEventHandler))
	c.Fire(&BaseEvent{target: "foo"})
	c.Fire(&BaseEvent{target: "bar"})
	c.Fire(&BaseEvent{target: "exit"})
	c.Run(1)
	if len(calls) != 2 {
		t.Errorf("Expected two calls to the EventHandler. Got %d.", len(calls))
	} else if calls[0] != "foo" || calls[1] != "bar" {
		t.Error("Expected event \"foo\" to occur before event \"bar\".")
	}
}

func Test_UnregisterEventHandler(t *testing.T) {
	calls = make([]string, 0)
	c := NewComponent()
	eventHandler := NewEventHandler("foo", BasicEventHandler)
	c.RegisterEventHandler(eventHandler)
	c.UnregisterEventHandler(eventHandler)
	c.Fire(&BaseEvent{target: "foo"})
	c.Fire(&BaseEvent{target: "exit"})
	c.Run(1)
	if len(calls) != 0 {
		t.Errorf("Expected no calls to the EventHandler. Got %d.", len(calls))
	}
}

func Test_RegisterComponent(t *testing.T) {
	calls = make([]string, 0)
	main := NewComponent()
	child := NewComponent()
	child.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	main.RegisterComponent(child)
	main.Fire(&BaseEvent{target: "foo"})
	main.Fire(&BaseEvent{target: "exit"})
	main.Run(1)
	if len(calls) != 1 {
		t.Errorf("Expected one call to the event handler. Got %d.", len(calls))
	}
}

func Test_UnregisterComponent(t *testing.T) {
	calls = make([]string, 0)
	main := NewComponent()
	child := NewComponent()
	child.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	main.RegisterComponent(child)
	main.UnregisterComponent(child)
	main.Fire(&BaseEvent{target: "foo"})
	main.Fire(&BaseEvent{target: "exit"})
	main.Run(1)
	if len(calls) != 0 {
		t.Errorf("Expected no calls to the event handler. Got %d.", len(calls))
	}
}

func Test_CompleteNotification(t *testing.T) {
	calls = make([]string, 0)
	c := NewComponent()
	eventHandler := NewEventHandler("foo", BasicEventHandler)
	c.RegisterEventHandler(eventHandler)
	completeHandler := NewEventHandler("foo_complete", BasicEventHandler)
	c.RegisterEventHandler(completeHandler)
	c.Fire(&BaseEvent{target: "foo", notify_complete: true})
	c.Fire(&BaseEvent{target: "exit"})
	c.Run(1)
	if len(calls) != 2 {
		t.Errorf("Expected two calls to the event handler. Got %d.", len(calls))
	}
}

func Test_SuccessNotification(t *testing.T) {
        calls = make([]string, 0)
        c := NewComponent()
        eventHandler := NewEventHandler("foo", BasicEventHandler)
        c.RegisterEventHandler(eventHandler)
        completeHandler := NewEventHandler("foo_success", BasicEventHandler)
        c.RegisterEventHandler(completeHandler)
        c.Fire(&BaseEvent{target: "foo", notify_success: true})
        c.Fire(&BaseEvent{target: "exit"})
        c.Run(1)
        if len(calls) != 2 {
                t.Errorf("Expected two calls to the event handler. Got %d.", len(calls))
        }
}

func Test_FailureNotification(t *testing.T) {
        calls = make([]string, 0)
        c := NewComponent()
        eventHandler := NewEventHandler("foo", FailingEventHandler)
        c.RegisterEventHandler(eventHandler)
        completeHandler := NewEventHandler("foo_failure", BasicEventHandler)
        c.RegisterEventHandler(completeHandler)
        c.Fire(&BaseEvent{target: "foo", notify_failure: true})
        c.Fire(&BaseEvent{target: "exit"})
        c.Run(1)
        if len(calls) != 2 {
                t.Errorf("Expected two calls to the event handler. Got %d.", len(calls))
        }
}

