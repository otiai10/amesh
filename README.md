amesh
==========

[![Go Reference](https://pkg.go.dev/badge/github.com/otiai10/amesh.svg)](https://pkg.go.dev/github.com/otiai10/amesh)
[![Actions Status](https://github.com/otiai10/amesh/workflows/Go/badge.svg)](https://github.com/otiai10/amesh/actions?query=workflow%3AGo)
[![Release](https://github.com/otiai10/amesh/actions/workflows/release.yml/badge.svg)](https://github.com/otiai10/amesh/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/otiai10/amesh/branch/main/graph/badge.svg?token=2KSMSazCPd)](https://codecov.io/gh/otiai10/amesh)
[![Go Report Card](https://goreportcard.com/badge/github.com/otiai10/amesh)](https://goreportcard.com/report/github.com/otiai10/amesh)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/otiai10/amesh?sort=semver)](https://pkg.go.dev/github.com/otiai10/amesh)

みんな大好き東京アメッシュ http://tokyo-ame.jwa.or.jp/
をCLIで表示

| iTerm | Sixel | default |
|:-----:|:-----:|:-------:|
| <img width="320px" src="https://user-images.githubusercontent.com/931554/39689648-8e8520b4-5212-11e8-87e2-b0bad05f530c.png"> | <img width="320px" src="https://user-images.githubusercontent.com/10111/39798686-7d505878-539c-11e8-8671-322f495824cb.png"> | <img width="320px" src="https://cloud.githubusercontent.com/assets/931554/11038037/5940e5be-8744-11e5-94d9-4b0bc7b2f55f.png"> |


# Install

```
go get -u -v github.com/otiai10/amesh
```

なんかローカルで `go get` もしたくないしバイナリも持ちたくない、というひとがいたので謎にDockerコンテナで表示させるようにしました。

```sh
docker run -e TERM_PROGRAM --rm otiai10/amesh
# たぶん、
# alias amesh='docker run -e TERM_PROGRAM --rm otiai10/amesh'
# したら幸せになれる
```

# Usage

```sh
amesh      # 降雨状況と地形と地名・県境を表示
amesh -a   # 直近30分をタイムラプスで表示
amesh -g=0 # 地形情報を非表示
amesh -b=0 # 地名・県境を非表示
amesh -p   # iTermを使っててもピクセル表示
```

# Slackで @amesh って言うとアメッシュの画像出すbot

<img width="40%" src="https://user-images.githubusercontent.com/931554/44345661-e5c65a00-a4ce-11e8-96a3-a024b8651183.png" >

詳しくは、 https://github.com/otiai10/amesh-bot
