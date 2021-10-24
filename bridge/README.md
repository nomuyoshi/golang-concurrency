# bridgeチャネル
次のようなチャネルのシーケンスから値を取得したいときにbridgeチャネルというパターンが使える。  
チャネルのチャネルを崩して単一のチャネルにするというパターン。

```
<-chan <-chan interface{}
```

要は、「チャネルをまとめたチャネル」　→ 「単一のチャネル」にするパターン

利用者側としては、チャネルのチャネルを扱わずに済み、受け取った値の処理のみに集中できる。