# Contextパッケージ
doneチャネルの代わりとして使える便利パッケージ。  
キャンセル通知や期限（デッドライン）、その他のリクエストスコープの値を伝達する手段が欲しく、context パッケージが作られた。  
【参考】
- https://go.dev/blog/context
- https://pkg.go.dev/context

ざっくり要点をまとめると
- キャンセル通知や期限（デッドライン）、その他のリクエストスコープの値を伝達する手段
- Contextの伝播は、第一引数に ctx という名前で渡す
- 必要に応じて WithCancel、WithDeadline、WithTimeout、WithValueを使用して子Context（派生）を作成する
- 親Contextがキャンセルされると派生したすべての子Contextもキャンセルされる

