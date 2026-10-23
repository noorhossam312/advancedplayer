package downloader

import (
	"context"

	"github.com/colespringer/waxtap/v3"
)

func DownloadFile(link string, path string) error {
	client, err := waxtap.New(waxtap.Options{})
	if err != nil {
		return err
	}
	_, err = client.Download(context.Background(), waxtap.Request{
		URL: link,
		ProcessSpec: waxtap.ProcessSpec{
			Transcode: &waxtap.TranscodeSpec{Format: waxtap.FormatVorbis},
			Output:    waxtap.ToFile(path),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
