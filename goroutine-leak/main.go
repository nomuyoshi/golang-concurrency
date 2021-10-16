package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan struct{}, strings <-chan string) <-chan struct{} {
		terminated := make(chan struct{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println("stringsから受信: ", s)
				case <-done:
					// doneチャネルが閉じられたらreturnしてゴルーチンを終了させる
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan struct{})
	// nil を渡しているので、doWork内では何も受信することができず、ブロックされ続ける
	// doneチャネル自体が無い or doneチャネルを閉じわすれたらdoWork内のゴルーチンは存在し続けてしまう
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	<-terminated
	fmt.Println("DONE")
}
