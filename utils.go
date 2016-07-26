package main

import (
	"sync"
)

func async_run(threads int, c *Component, wg *sync.WaitGroup) {
	defer wg.Done()
	c.Run(threads)
}

