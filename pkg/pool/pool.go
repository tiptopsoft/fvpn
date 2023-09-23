package main

import (
	"fmt"
	"time"
)

type MyObject struct {
	Name  string
	Value int
	Other string
}

var myPool = NewPool(func() any {
	fmt.Println("create new object")
	obj := new(MyObject)
	obj.Name = "new one"
	return obj
})

func worker(id int, ch chan *MyObject) {

	for {
		obj := myPool.Get().(*MyObject)
		obj.Value = id
		if id == 2 {
			obj.Other = "other......"
		}
		id++

		ch <- obj
		time.Sleep(1 * time.Second)
	}

}

func main() {
	numWorkers := 1
	ch := make(chan *MyObject, numWorkers)

	go worker(0, ch)
	go worker(0, ch)

	go func() {
		for {
			obj := <-ch
			fmt.Printf("Received object with value %d, name: %s, other: %s\n", obj.Value, obj.Name, obj.Other)
			//obj.Value = 0
			myPool.Put(obj)
		}

		//myPool.Put(obj)
	}()

	time.Sleep(1000 * time.Second)
}
