package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	<-orChannel(
		sig(60*time.Second),
		sig(20*time.Second),
		sig(10*time.Second),
		sig(3*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
	)
	fmt.Printf("Done after %v\n", time.Since(start))
}

// orChannel は複数のチャネルをまとめた orDone チャネルを返す
// まとめたチャネルのうち、どれか1つでも閉じられたら orDone チャネルも閉じる
// 例. チャネル6個(ch1, ch2, ... , ch6)をまとめる場合
// 1回目
// ・ orDone(1回目)チャネルを生成
// ・ ch1, ch2, ch3から受信待ち（閉じられ待ち）
// ・ ch4, ch5, ch6, orDone(1回目)を引数にして再度orChannelを呼び出す（再帰）→ orDoneチャネル(2回目)の閉じられ待ち
// 2回目
// ・ orDone(2回目)チャネルを生成
// ・ ch4. ch5, ch6 から受信待ち（閉じられ待ち）
// ・ orDone(1回目), orDone(2回目)を引数にして再度orChannelを呼び出す（再帰）→ orDoneチャネル(3回目)の閉じられ待ち
// 3回目
// ・ orDone(3回目)チャネルを生成
// ・ orDone(1回目), orDone(2回目)から受信待ち（閉じられ待ち）
// 最終的に
// → ch1 ~ ch6のどれかが閉じらる
// → どれかのorDoneが閉じられる
// → 全てのorDoneが閉じられる
// → 呼び出し元から見ると、返り値のorDoneチャネル(1回目)が閉じられる
func orChannel(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0: // チャネルが0個ならnil
		return nil
	case 1: // チャネルが1個ならまとめる必要がないので、そのチャネルを返す
		return channels[0]
	}

	// 複数のチャネルをまとめるチャネル
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			// 2つだけなら、どちらかのチャネルが閉じられるまで待機
			//
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			// 3つ以上あるならorチャネルを再帰呼び出し
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			// 残りの部分のみをorチャネルに渡す。
			// 上位のチャネルが閉じたら下位も終了するようにorDoneチャネルも渡す。
			case <-orChannel(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

// sig は指定した時間経過したら閉じられるチャネルを返す
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}
