amesh
==========

みんな大好き東京アメッシュ http://tokyo-ame.jwa.or.jp/
をCLIで表示

![](https://cloud.githubusercontent.com/assets/931554/11038037/5940e5be-8744-11e5-94d9-4b0bc7b2f55f.png)

# install

```
go get github.com/otiai10/amesh/amesh
```

# usage

```sh
amesh     #降雨状況のみ描画
amesh -g  #地形情報を描画
amesh -m  #地名情報を描画
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
