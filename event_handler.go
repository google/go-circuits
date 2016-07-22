package main

func NewEventHandler(target string, c func(Event)) *EventHandler {
	eh := EventHandler{Call: c, target: target}
	return &eh
}

type EventHandler struct {
	Call func(Event)
	target string
}

func (eh *EventHandler) Target() string {
	return eh.target
}

