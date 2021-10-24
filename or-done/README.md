# or-done チャネル
システムの完全に異なる部分から受け取ったチャネルを扱うときなど、select文を連続するループを書く必要が出てくる場合がある。  
ゴルーチンを1つ使って複雑なselect文をカプセル化するとスッキリかける。

↓な状況の場合、orDoneチャネル関数を作ってカプセル化すると良い。
```go
for val := range myChan {
  // 何かする
}

loop:
  for {
    select {
      case <-done:
        break loop
      case val, ok := <-myChan:
        if ok == false {
          return // またはforからbreak
        }
        // valに対してなにかする
    }
  }
```
