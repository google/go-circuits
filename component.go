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
	"container/list"
	"sync"
)

type ExitEvent struct {
	BaseEvent
	Exit func()
}

func NewComponent() *Component {
	return NewAdvancedComponent(100)
}

func NewAdvancedComponent(eventChanSize int) *Component {
	c := Component{
		eventHandlers: make(map[string]map[string]*list.List),
		eventChan:     make(chan Event, eventChanSize),
	}
	c.RegisterEventHandler(NewEventHandler("*", "exit", c.Exit))
	return &c
}

type Component struct {
	root          *Component
	parent        *Component
	children      []*Component
	eventHandlers map[string]map[string]*list.List
	eventChan     chan Event
	waitGroup     sync.WaitGroup
}

// Add an EventHandler to the system
func (c *Component) RegisterEventHandler(handler *EventHandler) {
	channel := handler.Channel()
	target := handler.Target()
	channel_handlers, ok := c.eventHandlers[channel]
	if !ok {
		channel_handlers = make(map[string]*list.List)
		c.eventHandlers[channel] = channel_handlers
	}
	handlers, ok := channel_handlers[target]
	if !ok {
		handlers = list.New()
		channel_handlers[target] = handlers
	}
	handlers.PushBack(handler)

	if c.parent != nil {
		c.parent.RegisterEventHandler(handler)
	}
}

// Remove an EventHandler from the system
func (c *Component) UnregisterEventHandler(handler *EventHandler) {
	channel := handler.Channel()
	target := handler.Target()
	handlers, ok := c.eventHandlers[channel][target]
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

	for _, channel_handlers := range component.eventHandlers {
		for _, handlers := range channel_handlers {
			for h := handlers.Front(); h != nil; h = h.Next() {
				c.RegisterEventHandler(h.Value.(*EventHandler))
			}
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
		for _, channel_handlers := range component.eventHandlers {
			for _, handlers := range channel_handlers {
				for h := handlers.Front(); h != nil; h = h.Next() {
					c.UnregisterEventHandler(h.Value.(*EventHandler))
				}
			}
		}
	}
}

// Add an event to the queue
func (c *Component) Fire(event Event) {
	c.eventChan <- event
}

// Trigger an exit on the next loop
func (c *Component) Exit(_ Event) {
	close(c.eventChan)
}

// Process events piped over the channel
func (c *Component) processEvents() {
	defer c.waitGroup.Done()
	for true {
		event, open := <-c.eventChan
		if !open {
			break
		}
		c.processEvent(event)
	}
}

// Process a single event
func (c *Component) processEvent(event Event) {
	channel := event.Channel()
	target := event.Target()
	handlers := c.getEventHandlers(channel, target)
	for _, h := range handlers {
		err := h.Call(event)

		// Send notifications of Event status
		if event.NotifyFailure() && err != nil {
			c.Fire(NewEvent(channel, target+"_failure"))
		} else if event.NotifySuccess() && err == nil {
			c.Fire(NewEvent(channel, target+"_success"))
		}
		if event.NotifyComplete() {
			c.Fire(NewEvent(channel, target+"_complete"))
		}
	}
}

// Determine event handlers for an event
func (c *Component) getEventHandlers(channel, target string) []*EventHandler {
	handlers := list.New()
	if channel == "*" {
		for _, channel_handlers := range c.eventHandlers {
			c.addEventHandlersFromChannel(handlers, channel_handlers, target)
		}
	} else {
		c.addEventHandlersFromChannel(handlers, c.eventHandlers[channel], target)
		c.addEventHandlersFromChannel(handlers, c.eventHandlers["*"], target)
	}

	ret := make([]*EventHandler, 0, handlers.Len())
	for h := handlers.Front(); h != nil; h = h.Next() {
		ret = append(ret, h.Value.(*EventHandler))
	}
	return ret
}

func (c *Component) addEventHandlersFromChannel(handlers *list.List, channel_handlers map[string]*list.List, target string) {
	if target == "*" {
		for _, target_handlers := range channel_handlers {
			if target_handlers != nil {
				handlers.PushBackList(target_handlers)
			}
		}
	} else {
		target_handlers := channel_handlers[target]
		if target_handlers != nil {
			handlers.PushBackList(target_handlers)
		}
		generic_handlers := channel_handlers["*"]
		if generic_handlers != nil {
			handlers.PushBackList(generic_handlers)
		}
	}
}

// Main loop
func (c *Component) Run(num_workers int) {
	if c.root != nil {
		panic("Cannot run main loop on child component.")
	}

	for i := 0; i < num_workers; i++ {
		c.waitGroup.Add(1)
		go c.processEvents()
	}

	c.waitGroup.Wait() // Wait for all event goroutines to finish
}
