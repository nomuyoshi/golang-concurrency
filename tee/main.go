package main

import (
	"fmt"
	"sync"
)

func main() {
	in := make(chan interface{})
	done := make(chan interface{})
	defer close(done)
	// inに値を詰める
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			select {
			case <-done:
				return
			case in <- i:
			}
		}
	}()

	// inに詰めた値をout1, out2の両方に詰める
	out1, out2 := tee(done, in)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range out1 {
			fmt.Println("out1の値: ", v)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range out2 {
			fmt.Println("out2の値: ", v)
		}
	}()

	wg.Wait()
	return
}

// teeチャネル 読み込み元のチャネルを渡し、同じ値を持つ2つのチャネルが返される。
func tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			// チャネルのコピー変数としてローカル変数を用意
			var out1, out2 = out1, out2
			// out1,2 に値を送信するために2回繰り返す
			for i := 0; i < 2; i++ {
				select {
				case out1 <- val:
					// 送信後にnilを入れてもう片方のチャネルに書き込まれるようにする
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}
