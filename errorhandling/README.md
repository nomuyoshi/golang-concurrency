# エラーハンドリング
同期関数を書くときを同じように考えれば良い。  
【ポイント】  
- エラーに対処するゴルーチンとエラーの生成者のゴルーチンを切り分ける(関心事を分ける)
- ゴルーチンがエラーを生成するのであれば、それは正常系の結果と同じ経路を使って渡されるべき。
- 渡された側でエラーに対して適切に対処する。

badとgoodのサンプルコード参照
