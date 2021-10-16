package main

import (
	"bytes"
	"fmt"
	"sync"
)

// 1と2によって予期せぬアクセスが起こらないようにしている
func main() {
	// 1. dataに直接アクセスせずに、引数として受け取っている
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	data := []byte("Hello, Golang.")
	wg.Add(2)
	// 2. 起動したゴルーチンがdataの一部しかアクセスできないように拘束
	go printData(&wg, data[:5])
	go printData(&wg, data[5:])
	wg.Wait()
	fmt.Println("終了")
}
