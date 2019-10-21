package main

import (
	"./threadpool"
)

type test1 struct {
}

func (this *test1) Run() bool {
	return true
}
func (this *test1) PsentMainThread() int {
	return threadpool.THREAD_COMPLETE
}

func main() {

	ThreadPool := threadpool.New()
	ThreadPool.AddTask(&test1{})
	ThreadPool.AddTask(&test1{})
	ThreadPool.AddTask(&test1{})
	ThreadPool.AddTask(&test1{})
	for {
		ThreadPool.Update()
	}
}
