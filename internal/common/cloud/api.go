package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

type CloudInfoProvider struct {
	client *http.Client
	mu     sync.RWMutex
	cache  map[string]ScyllaVersion
}

func NewCloudInfoProvider(ctx context.Context, client *http.Client) *CloudInfoProvider {
	api := &CloudInfoProvider{client: client}

	go api.LoadData(ctx)

	return api
}

type ScyllaVersion struct {
	ID          int    `json:"id"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NewCluster  string `json:"newCluster"`
}

type Data struct {
	ScyllaVersions []ScyllaVersion `json:"scyllaVersions"`
}

type Response struct {
	Data Data `json:"data"`
}

func (a *CloudInfoProvider) LoadData(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.cloud.scylladb.com/deployment/scylla-versions", nil)
	if err != nil {
		return
	}

	resp, err := a.client.Do(rq)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	a.cache = make(map[string]ScyllaVersion)
	for _, v := range data.Data.ScyllaVersions {
		a.cache[v.Version] = v
	}
}

const UnknownAvailbility = "N/A"
const ErrAvailability = "Error while fetching data"

func (a *CloudInfoProvider) CloudAvailability(version string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.cache == nil {
		return ErrAvailability
	}

	v, ok := a.cache[version]
	if !ok {
		return UnknownAvailbility
	}

	return v.NewCluster
}
