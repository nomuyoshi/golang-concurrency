package main

import (
	"fmt"
	"time"
)

// for-select はpipelineの例でも使っているので、そちらも参照
func main() {
	input := func() <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for _, s := range []string{"a", "b", "c", "d", "e", "f", "g"} {
				time.Sleep(1 * time.Millisecond)
				stringStream <- s
			}
		}()
		return stringStream
	}

	stringStream := input()
	// 無限ループまたはイテレーションを回す
	for {
		// チャネルに対して何か行う
		select {
		case s, ok := <-stringStream:
			if !ok {
				fmt.Println("DONE")
				return
			}
			fmt.Println("received value = ", s)
		}
	}
}
