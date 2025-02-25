package meilisearch

import (
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"net/http"
	"time"
)

type Meili struct {
	Client            meilisearch.ServiceManager
	DefaultCollection string
}

func New(cfg Config) *Meili {

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
	}
	meiliClient := meilisearch.New(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		meilisearch.WithCustomClient(httpClient))

	return &Meili{
		Client:            meiliClient,
		DefaultCollection: cfg.DefaultCollection,
	}
}
