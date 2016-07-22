package main

import (
	"container/list"
	"time"
)

func NewComponent() *Component {
	c := Component{eventHandlers: make(map[string]*list.List), eventQueue: list.New()}
	return &c
}

type Component struct {
	root *Component
	parent *Component
	children []*Component
	eventHandlers map[string]*list.List
	eventQueue *list.List
	exit bool
}

// Add an EventHandler to the system
func (c *Component) RegisterEventHandler(handler *EventHandler) {
	target := handler.Target()
	_, ok := c.eventHandlers[target]
	if !ok {
		c.eventHandlers[target] = list.New()
	}
	c.eventHandlers[target].PushBack(handler)

	if c.parent != nil {
		c.parent.RegisterEventHandler(handler)
	}
}

// Remove an EventHandler from the system
func (c *Component) UnregisterEventHandler(handler *EventHandler) {
	target := handler.Target()
	handlers, ok := c.eventHandlers[target]
	if !ok {
		return
	}

	for h := handlers.Front(); h != nil; h = h.Next() {
		if h.Value.(*EventHandler) == handler {
			handlers.Remove(h)
			break
		}
	}

	if c.parent != nil {
		c.parent.UnregisterEventHandler(handler)
	}
}

// Add a child Component to the system, implicitly adds all of children recursively
func (c *Component) RegisterComponent(component *Component) {
	c.children = append(c.children, component)
	component.root = c.root
	component.parent = c

	for _, handlers := range component.eventHandlers {
		for h := handlers.Front(); h != nil; h = h.Next() {
			c.RegisterEventHandler(h.Value.(*EventHandler))
		}
	}
}

// Remove a Component from the system, implicitly removes all children recursively
func (c *Component) UnregisterComponent(component *Component) {
	found := false
	i := 0
	for ; i < len(c.children); i++ {
		if c.children[i] == component {
			found = true
			break
		}
	}

	if found {
		c.children = append(c.children[:i], c.children[i+1:]...)
		for _, handlers := range component.eventHandlers {
			for h := handlers.Front(); h != nil; h = h.Next() {
				c.UnregisterEventHandler(h.Value.(*EventHandler))
			}
		}
	}
}

// Add an event to the queue
func (c *Component) Fire(event Event) {
	c.eventQueue.PushBack(event)
}

// Trigger an exit on the next loop
func (c *Component) Exit() {
	c.exit = true
}

// Process all events in the queue
func (c *Component) executeQueue() {
	for e := c.eventQueue.Front(); c.eventQueue.Len() != 0; e = c.eventQueue.Front() {
		c.eventQueue.Remove(e)
		event := e.Value.(Event)
		handlers := c.eventHandlers[event.GetTarget()]
		if handlers == nil {
			continue
		}
		for h := handlers.Front(); h != nil; h = h.Next() {
			h.Value.(*EventHandler).Call(event)
		}
	}
}

// Recursively tick over all children
func (c *Component) recursiveTick() {
	c.Tick()
	for _, child := range c.children {
		child.recursiveTick()
	}
}

// Iterate once
func (c *Component) Tick() {
	c.executeQueue()
}

// Main loop
func (c *Component) Run() {
	if c.root != nil {
		panic("Cannot run main loop on child component.")
	}

	for !c.exit {
		c.recursiveTick()
		time.Sleep(1 * time.Millisecond)
	}
}

