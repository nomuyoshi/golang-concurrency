package main

import (
	"fmt"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	generator := func(done <-chan interface{}, ints ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, v := range ints {
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// 入力チャネルと加算値を受け取り、加算値を加えた値を出力チャネルに送信
	add := func(done <-chan interface{}, in <-chan int, additive int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for v := range in {
				select {
				case out <- v + additive:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// 入力チャネルと乗数を受け取り、乗数を掛けた値を出力チャネルに送信
	multiply := func(done <-chan interface{}, in <-chan int, multiplier int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for v := range in {
				select {
				case out <- v * multiplier:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// generator でデータをチャネル（ストリーム）に変換
	intStream := generator(done, 1, 2, 3, 4)
	// pipeline処理を構築
	pipeline := multiply(done, add(done, intStream, 10), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}
