package controllers

import (
	"bytes"
	"encoding/json"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"google.golang.org/appengine/urlfetch"

	"github.com/otiai10/amesh"
	"github.com/otiai10/marmoset"
	"google.golang.org/appengine"
)

// Webhook ...
func Webhook(w http.ResponseWriter, r *http.Request) {
	render := marmoset.Render(w, true)

	entry := amesh.GetEntry()

	// {{{ TODO: AppEngine特有のコードなのでビルドフラグ付きミドルウェアにする
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	// }}}

	img, err := entry.Image(true, true, client)
	if err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	// TODO: Separate to platforms
	// {{{ Create multipart form data and files here
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	f, err := writer.CreateFormFile("file", "amesh.png")
	if err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	if _, err := io.Copy(f, buf); err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	if err := writer.WriteField("token", os.Getenv("SLACK_BOT_ACCESS_TOKEN")); err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	if err := writer.WriteField("channels", os.Getenv("SLACK_CHANNELS")); err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	if err := writer.Close(); err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", body)
	if err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	if res.StatusCode != http.StatusOK {
		render.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": res.Status,
		})
		return
	}
	defer res.Body.Close()
	// }}}

	response := map[string]interface{}{}
	json.NewDecoder(res.Body).Decode(&response)

	render.JSON(http.StatusOK, map[string]interface{}{
		"status":   res.StatusCode,
		"message":  res.Status,
		"response": response,
	})
}
