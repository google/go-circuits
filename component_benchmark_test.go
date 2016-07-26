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
	"runtime"
	"sync"
        "testing"
)

func benchmark(b *testing.B, threads int) {
        c := NewComponent()
        c.RegisterEventHandler(NewEventHandler("f", func(_ Event){}))
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go async_run(threads, c, wg)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                c.Fire(BaseEvent{"f"})
        }
	c.Fire(BaseEvent{"exit"})
        wg.Wait()
}

func Benchmark_SingleThread(b *testing.B) {
	benchmark(b, 1)
}

func Benchmark_MultiThread(b *testing.B) {
	benchmark(b, runtime.GOMAXPROCS(0))
}
