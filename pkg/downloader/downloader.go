package downloader

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpDownloader struct {
	Timeout int
}

func (d *HttpDownloader) Download(url string) (string, error) {
	client := &http.Client{
		Timeout: time.Duration(d.Timeout) * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("url %s responded with status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
