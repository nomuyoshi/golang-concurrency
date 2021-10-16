package main

import "fmt"

func main() {
	chanOwner := func() <-chan int {
		// チャネルをchanOwner関数のレキシカルスコープ内で初期化
		// -> chチャネルへ書き込みできるスコープを制限
		// -> チャネルへの書き込み権限を拘束して、他のゴルーチンによる書き込みを防いでいる
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; i < 5; i++ {
				ch <- i
			}
		}()
		return ch
	}

	consumer := func(ch <-chan int) {
		for v := range ch {
			fmt.Println("受信値: ", v)
		}
		fmt.Println("受信操作終了")
	}

	ch := chanOwner()
	consumer(ch)
}
