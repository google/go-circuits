package main

import (
	//"fmt"
	"testing"
)

var calls []string

func BasicEventHandler(e Event) {
	calls = append(calls, e.GetTarget())
}

func Test_SimpleEvent(t *testing.T) {
	calls = make([]string, 0)
	c := NewComponent()
	c.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	c.Fire(BaseEvent{"foo"})
	c.Tick()
	if len(calls) != 1 {
		t.Errorf("Expected one call to the EventHandler. Got %d.", calls)
	}
}

func Test_EventsFIFO(t *testing.T) {
	calls = make([]string, 0)
	c := NewComponent()
	c.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	c.RegisterEventHandler(NewEventHandler("bar", BasicEventHandler))
	c.Fire(BaseEvent{"foo"})
	c.Fire(BaseEvent{"bar"})
	c.Tick()
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
	c.Fire(BaseEvent{"foo"})
	c.Tick()
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
	main.Fire(BaseEvent{"foo"})
	main.Tick()
	if len(calls) != 1 {
		t.Errorf("Expected one call to the event handler. Got %d.", len(calls))
	}
}

func Test_UnregisterComponent(t *testing.T) {
	calls := make([]string, 0)
	main := NewComponent()
	child := NewComponent()
	child.RegisterEventHandler(NewEventHandler("foo", BasicEventHandler))
	main.RegisterComponent(child)
	main.Fire(BaseEvent{"foo"})
	main.Tick()
	if len(calls) != 0 {
		t.Errorf("Expected no calls to the event handler. Got %d.", len(calls))
	}
}
