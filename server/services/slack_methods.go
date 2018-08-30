package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/otiai10/amesh/plugins/typhoon"
	"github.com/otiai10/chant/server/middleware/lib/google"

	"github.com/otiai10/amesh"
	"github.com/otiai10/amesh/server/middlewares"
	m "github.com/otiai10/marmoset"
)

func (slack *Slack) methodClean(ctx context.Context, channel, bot string) error {
	files, err := slack.listUploadedFiles(ctx, channel, bot)
	if err != nil {
		return err
	}
	if _, err := slack.postMessage(ctx, fmt.Sprintf("%d files found.", len(files)), channel); err != nil {
		return err
	}
	if err := slack.deleteUploadedFiles(ctx, files); err != nil {
		return err
	}
	if _, err := slack.postMessage(ctx, "cleaned up", channel); err != nil {
		return err
	}
	return nil
}

func (slack *Slack) listUploadedFiles(ctx context.Context, channel, bot string) ([]*SlackFile, error) {
	query := url.Values{}
	query.Add("token", slack.UserAccessToken)
	query.Add("channel", channel)
	query.Add("user", bot)
	query.Add("types", "images")
	res, err := middlewares.HTTPClient(ctx).Get("https://slack.com/api/files.list" + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	response := new(SlackAPIResponse)
	json.NewDecoder(res.Body).Decode(response)
	if !response.OK {
		return nil, fmt.Errorf("SLACK API files.list: %s", response.Error)
	}
	return response.Files, nil
}

func (slack *Slack) deleteUploadedFiles(ctx context.Context, files []*SlackFile) error {
	for _, f := range files {
		if err := slack.deleteUploadedFile(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
func (slack *Slack) deleteUploadedFile(ctx context.Context, file *SlackFile) error {
	if file == nil {
		return nil
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(m.P{"file": file.ID})
	req, err := http.NewRequest("POST", "https://slack.com/api/files.delete", body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", slack.BotAccessToken))
	res, err := middlewares.HTTPClient(ctx).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	response := new(SlackAPIResponse)
	json.NewDecoder(res.Body).Decode(response)
	if !response.OK {
		return fmt.Errorf("SLACK API files.delete: %v", response.Error)
	}
	return nil
}

func (slack *Slack) uploadFile(ctx context.Context, file io.Reader, channel string) error {

	postbody := new(bytes.Buffer)
	writer := multipart.NewWriter(postbody)

	f, err := writer.CreateFormFile("file", "amesh.png")
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, file); err != nil {
		return err
	}

	if err := writer.WriteField("token", slack.BotAccessToken); err != nil {
		return err
	}

	if err := writer.WriteField("channels", channel); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", postbody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := middlewares.HTTPClient(ctx).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	response := new(SlackAPIResponse)
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}
	if !response.OK {
		return fmt.Errorf("SLACK API files.upload: %s", response.Error)
	}

	return nil
}

func (slack *Slack) methodTyphoon(ctx context.Context, channel string) error {

	client := middlewares.HTTPClient(ctx)

	entry, err := typhoon.GetEntry(client)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?%s=%d", entry.NearJP, "t", time.Now().Unix())
	_, err = slack.postMessage(ctx, url, channel)
	return err
}

func (slack *Slack) methodImageSearch(ctx context.Context, channel, query string) error {

	middlewares.Log(ctx).Debugf("methodImageSearch: %+v\n", query)
	// TODO: ちょっとめんどくさいんで otiai10/chant/middleware/lib/google 呼んでますけど
	//       これどっかにpackage分離しましょうねｗ
	client, err := google.NewClient(ctx)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())
	res, err := client.SearchImage(query, rand.Intn(10)+1)

	if err != nil {
		return err
	}
	if len(res.Items) == 0 {
		_, err := slack.postMessage(ctx, "ない", channel)
		return err
	}
	_, err = slack.postMessage(ctx, res.RandomItem().Link, channel)

	return err
}

func (slack *Slack) methodShow(ctx context.Context, channel string) error {

	client := middlewares.HTTPClient(ctx)

	entry := amesh.GetEntry()
	img, err := entry.Image(true, true, client)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return err
	}

	return slack.uploadFile(ctx, buf, channel)
}
