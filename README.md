amesh
==========

みんな大好き東京アメッシュ http://tokyo-ame.jwa.or.jp/
をCLIで表示

**iTermの場合**

（すごい雨降ってるときの画像）

![](https://user-images.githubusercontent.com/931554/39689648-8e8520b4-5212-11e8-87e2-b0bad05f530c.png)

**Sixel拡張をサポートするターミナルの場合**

（晴れてるときの画像）

![](https://user-images.githubusercontent.com/10111/39798686-7d505878-539c-11e8-8671-322f495824cb.png)

**それ意外のターミナルアプリ**

（千葉のほうだけちょっと雨降ってるときの画像）

![](https://cloud.githubusercontent.com/assets/931554/11038037/5940e5be-8744-11e5-94d9-4b0bc7b2f55f.png)

# install

```
go get github.com/otiai10/amesh/amesh
```

なんかローカルで `go get` もしたくないしバイナリも持ちたくない、というひとがいたので謎にDockerコンテナで表示させるようにしました。

```sh
docker run -e TERM_PROGRAM --rm otiai10/amesh
# たぶん、
# alias amesh='docker run -e TERM_PROGRAM --rm otiai10/amesh'
# したら幸せになれる
```

# usage

```sh
amesh      #降雨状況と地形と地名・県境を表示
amesh -g=0 #地形情報を非表示
amesh -b=0 #地名・県境を非表示
amesh -p   #iTermを使っててもピクセル表示
```

# daemon

3分おきにアメッシュ見に行ってSlackとかTwitterとかで通知してくれるデーモンが立ちます

```sh
touch develop.env
docker-compose up
```

# package

雨天判定とか、その他もろもろカスタマイズしたいときは、Observerをいじれます

```go
observer := amesh.NewObserver()

// 雨天判定をするためのメソッド
// デフォルトでは「雨天ピクセルが全体の30%以上」という判定してます
observer.IsRaining = func(ev amesh.Event) bool {
  // EventにはImgというimage.Imageが含まれてるので
  // お好きなアルゴリズムで雨天判定すればいいと思う
  r, g, b, a := ev.Img.At(320, 200).RGBA()
  return r*g*b*a != 0
}

observer.On(Rain, func(ev amesh.Event) error {
    // IsRain == true だったときに呼ばれるハンドラ
    return nil
})
observer.On(Update, func(ev amesh.Event) error {
    // Rain以外のクロールで必ず呼ばれるハンドラ
    if ev.Timestamp.Hour() == 23 {
      return fmt.Errorf("なにかしらエラー")
    }
    return nil
})
observer.On(Error, func(ev amesh.Event) error {
    // なんかエラーがあったときに呼ばれるハンドラ
    observer.Start() // リカバーして再起動したり
})

```
