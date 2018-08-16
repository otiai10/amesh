# Slackとかで /amesh で表示できるようにしたいというプロジェクト

`app/secret.yaml` に以下の内容が必要です

```yaml
env_variables:
    # アメッシュ画像の投稿先にSlackを使用
    SERVICE: Slack
    # SLackを使用する上で必要な変数
    SLACK_BOT_ACCESS_TOKEN: xoxb-12345-67890-xxxxxxxxxxx
    SLACK_VERIFICATION: AbcDefGhiJklMnoPqrStuVwxYz
```

そんで、

```sh
% goapp deploy server/app
```

# 開発

```sh
% goapp serve server/app
```
