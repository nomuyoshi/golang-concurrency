package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 1st stage: gen 複数のintを受け取り、出力チャネルに送信
// 2nd stage: sq 入力チャネルを受け取り、チャネルから値を受信して、二乗した値を出力チャネルに送信
// 3rd stage: fan-in 複数の2nd stageの出力チャネルを1つのチャネルにまとめる
// final stage: 入力チャネルから値を受信して出力

// main はパイプラインのセットアップと最終ステージの実行
func main() {
	done := make(chan struct{})
	// 実際には、deferでcloseする方が安全
	// defer close(done)

	in := gen(done, 1, 2, 3, 4, 5)
	c1 := sq(done, in)
	c2 := sq(done, in)

	out := merge(done, c1, c2)
	// 1つだけ受信して、あとの値は捨てる
	fmt.Println(<-out)
	// doneチャネルをクローズして、開始しているgoroutineを終了させる
	close(done)
	time.Sleep(2 * time.Second)
	fmt.Println("goroutine数 = ", runtime.NumGoroutine())
	fmt.Println("終了")
}

// gen は1st stage
// 複数のintを受け取り、出力チャネルに送信
// doneチャネルがクローズされたら、処理を中断する
func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				fmt.Println("gen canceled")
				return
			}
		}
		fmt.Println("gen finished")
	}()

	return out
}

// sq は 2nd stage
// 入力チャネルから値を受信し、二乗して出力チャネルに送信
// doneチャネルがクローズされたら、処理を中断する
func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				fmt.Println("sq canceled")
				return
			}
		}
		fmt.Println("sq finished")
	}()

	return out
}

// merge は 3rd stage (fan-in)
// 複数の入力チャネルを受け取り、1つの出力チャネルにまとめる
// doneチャネルがクローズされたら、処理を中断する
func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(c <-chan int) {
		defer wg.Done()

		for n := range c {
			select {
			case out <- n:
			case <-done:
				fmt.Println("output canceled")
				return
			}
		}
		fmt.Println("output finished")
	}

	for _, c := range cs {
		wg.Add(1)
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
		fmt.Println("wg.Wait finished")
	}()

	return out
}
