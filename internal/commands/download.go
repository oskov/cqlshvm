package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oskov/cqlshvm/internal/common/downloader"
	"github.com/oskov/cqlshvm/internal/common/reader"
	"github.com/oskov/cqlshvm/internal/common/version"
)

type DownloadCommand struct {
	Downloader *downloader.Downloader
	Reader     *reader.Reader
}

func NewDownloadComand(client *http.Client) *DownloadCommand {
	return &DownloadCommand{
		Downloader: downloader.NewDownloader(client),
		Reader:     reader.NewReader(client),
	}
}

func (c *DownloadCommand) Run(ctx context.Context, v version.Version) error {
	versions, err := c.Reader.ReadFilesVersions(ctx, reader.ListParams{
		Exact: &v,
	})
	if err != nil {
		return fmt.Errorf("failed to read versions: %w", err)
	}
	if len(versions) == 0 {
		return fmt.Errorf("version %s not found", v)
	}

	return c.Downloader.DownloadFile(versions[0].Key)
}
