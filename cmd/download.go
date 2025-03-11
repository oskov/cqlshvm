package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/oskov/cqlshvm/internal/commands"
	"github.com/oskov/cqlshvm/internal/common/version"

	"github.com/urfave/cli/v2"
)

var downloadCommand = &cli.Command{
	Name:      "download",
	Usage:     "Fetches the archive of the specified version (chosen from the output of the list command) and outputs the downloaded file to STDOUT.",
	UsageText: "cqlshvm download <version> > cqlsh.tar.gz",
	Args:      true,
	Action: func(cCtx *cli.Context) error {
		v := cCtx.Args().Get(0)

		parsed, err := version.Parse(v)
		if err != nil {
			return fmt.Errorf("failed to parse version: %w", err)
		}

		httpClient := &http.Client{
			Timeout: 10 * time.Second,
		}

		ctx, cl := context.WithTimeout(cCtx.Context, time.Second*30)
		defer cl()

		downloadCommand := commands.NewDownloadComand(httpClient)

		downloadCommand.Run(ctx, parsed)

		os.Stderr.WriteString("Download completed\n")

		return nil
	},
}
