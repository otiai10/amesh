package plugins

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/otiai10/chant/server/middleware/lib/google"
)

// Image ...
type Image struct {
	GoogleAPIKey               string
	GoogleCustomSearchEngineID string
}

// Method ...
func (search Image) Method() string {
	return "img"
}

// Match ...
func (search Image) Match(ctx context.Context, texts []string) bool {
	if len(texts) == 0 {
		return false
	}
	return texts[0] == "img"
}

// TaskValues ...
func (search Image) TaskValues(ctx context.Context, texts []string) url.Values {
	return url.Values{"query": {strings.Join(texts[1:], "+")}}
}

// Exec ...
// TODO: たぶんstringじゃないほうがいいんだよねえ
func (search Image) Exec(ctx context.Context, r *http.Request) (string, error) {

	query := r.FormValue("query")

	// TODO: ちょっとめんどくさいんで otiai10/chant/middleware/lib/google 呼んでますけど
	//       これどっかにpackage分離しましょうねｗ
	// TODO: ここで環境変数渡すのダサすぎるので、google.Clientの初期化手順に修正が必要
	os.Setenv("GOOGLE_SEARCH_API_KEY", search.GoogleAPIKey)
	os.Setenv("GOOGLE_SEARCH_ENGINE_ID", search.GoogleCustomSearchEngineID)
	client, err := google.NewClient(ctx)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().Unix())
	res, err := client.SearchImage(query, rand.Intn(10)+1)

	if err != nil {
		return "", err
	}

	if len(res.Items) == 0 {
		return "", fmt.Errorf("ないです")
	}

	return res.RandomItem().Link, nil
}
