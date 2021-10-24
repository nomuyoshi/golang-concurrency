package main

import "fmt"

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	done := make(chan interface{})
	defer close(done)
	// bridgeチャネルのおかげで「チャネルのチャネル」を1つのrangeループで処理できる
	for v := range bridge(done, genVals()) {
		fmt.Println(v)
	}
}

func bridge(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			// チャネルのチャネル(chanStream)からチャネルを受信して、streamに入れる
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			// streamから値を受信して、valStreamに詰める
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
			// streamが閉じられたら、次のループ（次のチャネルのチャネルに移る）
		}
	}()
	// valStream にはチャネルのチャネル(chanStream)に送られた全ての値が送信されている
	return valStream
}

// orDone は複雑なfor-selectをラップした関数
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
