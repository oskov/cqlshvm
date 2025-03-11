package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ryanuber/columnize"

	"github.com/oskov/cqlshvm/internal/commands"
	"github.com/oskov/cqlshvm/internal/common/version"

	"github.com/urfave/cli/v2"
)

var listCommand = &cli.Command{
	Name:      "list",
	Usage:     "Show all available versions of cqlsh and availability in ScyllaDB Cloud",
	ArgsUsage: "-lt <version> -gt <version>",
	UsageText: "cqlshvm list [-lt <version>] [-gt <version>]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "lt",
			Usage: "List only versions less than the specified version",
		},
		&cli.StringFlag{
			Name:  "gt",
			Usage: "List only versions greater than the specified version",
		},
	},
	Action: func(cCtx *cli.Context) error {
		lt := cCtx.String("lt")
		gt := cCtx.String("gt")

		httpClient := &http.Client{
			Timeout: 10 * time.Second,
		}
		ctx, cl := context.WithTimeout(cCtx.Context, time.Second*30)
		defer cl()

		listCommand := commands.NewListCommand(ctx, httpClient)

		params := commands.ListParams{}
		if lt != "" {
			v, err := version.Parse(lt)
			if err != nil {
				return err
			}
			params.Lt = &v
		}

		if gt != "" {
			v, err := version.Parse(gt)
			if err != nil {
				return err
			}
			params.Gt = &v
		}

		result, err := listCommand.Run(ctx, params)
		if err != nil {
			return err
		}

		var lines []string = []string{"VERSION|CLOUD"}

		for _, v := range result.Versions {
			lines = append(lines, fmt.Sprintf("%s|%s", v.Name, v.CloudAvailablity))
		}

		resultString := columnize.SimpleFormat(lines)
		fmt.Println(resultString)

		return nil
	},
}
