package commands

import (
	"context"
	"net/http"
	"slices"

	"github.com/oskov/cqlshvm/internal/common/cloud"
	"github.com/oskov/cqlshvm/internal/common/reader"
	"github.com/oskov/cqlshvm/internal/common/version"
)

type ListCommand struct {
	cloudInfoProvider *cloud.CloudInfoProvider
	reader            *reader.Reader
}

func NewListCommand(ctx context.Context, client *http.Client) *ListCommand {
	return &ListCommand{
		cloudInfoProvider: cloud.NewCloudInfoProvider(ctx, client),
		reader:            reader.NewReader(client),
	}
}

type ListVersion struct {
	Name             string
	CloudAvailablity string
}

type ListResult struct {
	Versions []ListVersion
}

type ListParams struct {
	Gt *version.Version
	Lt *version.Version
}

func (c *ListCommand) Run(ctx context.Context, params ListParams) (ListResult, error) {
	var result ListResult
	versions, err := c.reader.ReadFilesVersions(ctx, reader.ListParams{
		Gt: params.Gt,
		Lt: params.Lt,
	})
	if err != nil {
		return result, err
	}

	sortFn := func(a, b reader.ReadResult) int {
		if a.Version.Eq(b.Version) {
			return 0
		}
		if a.Version.Gt(b.Version) {
			return 1
		}
		return -1
	}

	slices.SortFunc(versions, sortFn)

	for _, v := range versions {
		result.Versions = append(result.Versions, ListVersion{
			Name:             v.Version.String(),
			CloudAvailablity: c.cloudInfoProvider.CloudAvailability(v.Version.String()),
		})
	}

	return result, nil
}
