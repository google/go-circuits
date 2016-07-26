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

import "errors"

func NewEventHandler(target string, c func(Event)) *EventHandler {
	eh := EventHandler{call: c, target: target}
	return &eh
}

type EventHandler struct {
	call func(Event)
	target string
}

func (eh *EventHandler) Target() string {
	return eh.target
}

func (eh *EventHandler) Call(event Event) (err error) {
	defer func() {
		if err_s := recover(); err_s != nil {
			err = errors.New(err_s.(string))
		}
	}()
	eh.call(event)
	return
}

