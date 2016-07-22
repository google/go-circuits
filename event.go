package main

type Event interface {
	GetTarget() string
}


type BaseEvent struct {
	target string
}

func (e BaseEvent) GetTarget() string {
	return e.target
}

