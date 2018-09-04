# Slackとかで @amesh って言うとアメッシュの画像出すやつ

<img width="40%" src="https://user-images.githubusercontent.com/931554/44345661-e5c65a00-a4ce-11e8-96a3-a024b8651183.png" >

# サービス側での設定

## Slack

1. アプリの作成
2. Botユーザの追加
3. Events API で subscribe　に `app_mention` を追加
    - エンドポイントに、以下のGAEの `/webhook/slack` を追加
4. Botのアクセストークンと、webhookのベリフィケーショントークンをコピる ( `secret.yaml` で使います )

# 必要な変数

`app/secret.yaml` に以下の内容が必要です

```yaml
env_variables:

    # アメッシュ画像の投稿先にSlackを使用
    SERVICE: Slack

    # SLackを使用する上で必要な変数
    SLACK_BOT_ACCESS_TOKEN: xoxb-12345-67890-xxxxxxxxxxx
    SLACK_VERIFICATION: AbcDefGhiJklMnoPqrStuVwxYz

    # アップロード枠の掃除のための `@amesh clean` を打つために必要
    SLACK_USER_ACCESS_TOKEN: xoxp-12345-67890-xxxxxxxxxx
```

# 必要な権限

チャットサービスとしてSlackを使う場合、作成したアプリは以下のScopesが必要になります

- `bot`: botユーザの追加や、botとしての発言を可能にするscope
- `files:read`: botユーザが投稿したアメッシュ画像の管理に必要なscope
- `channels:history`: botユーザが投稿した発言を取り消す場合に必要なscope

# サーバのデプロイ

```sh
% goapp deploy server/app
```

# 開発

```sh
% goapp serve server/app
```
