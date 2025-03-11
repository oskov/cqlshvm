package reader

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/url"
	"sync"

	"github.com/oskov/cqlshvm/internal/common/version"
)

var s3BucketURL = "https://s3.amazonaws.com/downloads.scylladb.com"

var defaultPrefix = "downloads/scylla-enterprise/relocatable/"

type Reader struct {
	client *http.Client
}

func NewReader(client *http.Client) *Reader {
	return &Reader{client: client}
}

type ListParams struct {
	Gt    *version.Version
	Lt    *version.Version
	Exact *version.Version
}

type ReadFilesVersionsResult struct {
	Name     string `xml:"Name"`
	Prefix   string `xml:"Prefix"`
	Contents []struct {
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         int    `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
	} `xml:"Contents"`
	CommonPrefixes []struct {
		Prefix string `xml:"Prefix"`
	} `xml:"CommonPrefixes"`
}

type ReadResult struct {
	Key     string
	Version version.Version
}

func (r *Reader) ReadFilesVersions(ctx context.Context, params ListParams) ([]ReadResult, error) {
	var result []ReadResult
	bucketObjects, err := r.gePage(ctx, defaultPrefix)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	results := make(chan []ReadResult, len(bucketObjects.CommonPrefixes))

	for _, v := range bucketObjects.CommonPrefixes {
		if !isValidPrefix(v.Prefix, params) {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()

			results <- r.processPrefixPage(ctx, params, v.Prefix)
		}()

	}

	wg.Wait()
	close(results)

	for v := range results {
		result = append(result, v...)
	}

	return result, nil
}

func (r *Reader) processPrefixPage(ctx context.Context, params ListParams, prefix string) []ReadResult {
	bucketObjects, err := r.gePage(ctx, prefix)
	if err != nil {
		return nil
	}

	result := []ReadResult{}

	for _, v := range bucketObjects.Contents {
		parsed, err := ParseObjectKey(v.Key)
		if err != nil {
			continue
		}
		if params.Gt != nil && !parsed.ObjectVersion.Gt(*params.Gt) {
			continue
		}
		if params.Lt != nil && !parsed.ObjectVersion.Lt(*params.Lt) {
			continue
		}
		if params.Exact != nil && !parsed.ObjectVersion.Eq(*params.Exact) {
			continue
		}

		result = append(result, ReadResult{
			Key:     v.Key,
			Version: parsed.ObjectVersion,
		})
	}

	return result
}

func prepareURL(bucketURL string, prefix string) (string, error) {
	u, err := url.Parse(bucketURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add("delimiter", "/")
	q.Add("prefix", prefix)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func isValidPrefix(prefix string, params ListParams) bool {
	parsed, err := ParseObjectKey(prefix)
	if err != nil {
		return false
	}

	if params.Gt != nil && params.Gt.Major > parsed.PrefixVersion.Major {
		return false
	}

	if params.Lt != nil && params.Lt.Major < parsed.PrefixVersion.Major {
		return false
	}
	if params.Exact != nil && params.Exact.Major != parsed.PrefixVersion.Major {
		return false
	}

	return parsed.PrefixVersion.Major >= 2024 // no point to check earlier years
}

func (r *Reader) gePage(ctx context.Context, prefix string) (ReadFilesVersionsResult, error) {
	var bucketObjects ReadFilesVersionsResult
	u, err := prepareURL(s3BucketURL, prefix)
	if err != nil {
		return bucketObjects, err
	}

	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return bucketObjects, err
	}

	resp, err := r.client.Do(rq)
	if err != nil {
		return bucketObjects, err
	}

	defer resp.Body.Close()
	err = xml.NewDecoder(resp.Body).Decode(&bucketObjects)
	if err != nil {
		return bucketObjects, err
	}

	return bucketObjects, nil
}
