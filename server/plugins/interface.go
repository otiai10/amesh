package plugins

import (
	"context"
	"net/http"
	"net/url"
)

// Plugin 直接アメッシュとは関係無いコマンドの実装インターフェース
type Plugin interface {

	// Match は、incomingな発言に対してこのプラグインが
	// 発火すぺきかどうかを真偽値で返す。
	Match(context.Context, []string) bool

	// TaskValues は、このプラグインが発火した結果、TaskQueueの
	// ワーカーにわたすパラメータを作成する。
	TaskValues(context.Context, []string) url.Values

	// Method は、GAE/TaskQueueのワーカーに、
	// このプラグインのタスクであることを伝えるためのキー。
	Method() string

	// Exec は、GAE/TaskQueueのワーカーで実際に行われる
	// 処理の具体的な内容を記述する。
	// TODO: stringを返すのだと、ちょっと表現力が足らん
	Exec(context.Context, *http.Request) (string, error)
}
