package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

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
