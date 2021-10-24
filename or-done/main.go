package main

import "fmt"

func main() {
	done := make(chan interface{})
	// myChanに値を送信する処理
	myChan := make(chan interface{})
	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
}

func orDone(done, ch <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-ch:
				if ok == false {
					return
				}
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()

	return valStream
}
