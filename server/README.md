amesh server
===============

usage

```sh
go run server/main.go
```

```sh
curl ":4010?pretty=true"
```

```json
{
	"url": "http://tokyo-ame.jwa.or.jp",
	"map": "http://tokyo-ame.jwa.or.jp/map/map000.jpg",
	"mesh": "http://tokyo-ame.jwa.or.jp/mesh/000/201510051930.gif",
	"mask": "http://tokyo-ame.jwa.or.jp/map/msk000.png"
}
```
