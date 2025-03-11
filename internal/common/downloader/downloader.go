package downloader

import (
	"io"
	"net/http"
	"os"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader(client *http.Client) *Downloader {
	return &Downloader{client: client}
}

const host = "https://s3.amazonaws.com/downloads.scylladb.com/"

func (d *Downloader) DownloadFile(key string) error {
	resp, err := d.client.Get(host + key)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)

	return err
}
