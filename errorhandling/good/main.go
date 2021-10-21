package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Response *http.Response
	Error    error
}

// 【ポイント】
// エラーに対処するゴルーチンとエラーの生成者のゴルーチンを切り分ける(関心事を分ける)
// ゴルーチンがエラーを生成するのであれば、それは正常系の結果と同じ経路を使って渡されるべき。
// そして、渡された側でエラーに対して適切に対処する。
// 同期関数を書くときを同じように考えれば良い。

// 今回の例だと
// ・checkStatusはレスポンスととエラーが対の「結果」を生成するだけ
// ・メインゴルーチンでは結果をみてエラーハンドリングする
// => エラーに対してメインゴルーチンで懸命な判断を下せる様になる
func main() {
	done := make(chan interface{})
	defer close(done)
	urls := []string{
		"https://www.google.com",
		"https://badhost",
		"https://www.yahoo.co.jp",
		"https://badhost",
		"https://qiita.com",
		"https://badhost",
		"https://twitter.com",
		"https://www.facebook.com",
	}
	var errCount int
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			errCount++
			fmt.Printf("error: %v\n", result.Error)
			// エラーが3つになったら処理を中断させる
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking.")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)
		for _, url := range urls {
			resp, err := http.Get(url)
			r := Result{Response: resp, Error: err}
			select {
			case <-done:
				return
			case results <- r:
			}
		}
	}()
	return results
}
