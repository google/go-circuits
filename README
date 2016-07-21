An Event driven Component framework for Golang based on Circuits
(http://github.com/circuits/circuits).

DISCLAIMER: This is not an official Google product.

Basic Usage

import (
    "fmt"
    "go-circuits"
)

type MyComponent struct {
    go-circuits.Component
}

func (mc *MyComponent) HelloWorldEventHandler(_ go-circuits.Event) {
    fmt.Println("Hello World!")
}

func NewMyComponent() *MyComponent {
    mc := MyComponent{}
    mc.RegisterEventHandler(go-circuits.NewEventHandler(
            "hello_world",
            mc.HelloWorldEventHandler
    ))
    return &mc
}

func main() {
    mc := NewMyComponent()
    mc.Fire(go-circuits.BaseEvent{"hello_world"})
    mc.Run()
}

