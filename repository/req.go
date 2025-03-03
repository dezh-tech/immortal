package repository

import (
	"context"
	"encoding/json"

	"github.com/dezh-tech/immortal/types"

	"github.com/meilisearch/meilisearch-go"

	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
)

func (h *Handler) HandleReq(f *filter.Filter, pubkey string) ([]event.Event, error) {
	meiliFilter := buildMeiliQuery(f)

	finalLimit := f.Limit
	if f.Limit <= 0 || f.Limit >= h.config.MaxQueryLimit {
		finalLimit = h.config.DefaultQueryLimit
	}

	sortBy := []string{"created_at:desc", "id:asc"}

	defaultCollection := h.meili.DefaultCollection

	searchResult, err := h.meili.Client.Index(defaultCollection).Search(f.Search,
		&meilisearch.SearchRequest{
			AttributesToSearchOn: []string{"content"},
			Limit:                int64(finalLimit),
			Sort:                 sortBy,
			Filter:               meiliFilter,
		})
	if err != nil {
		_, err := h.grpc.AddLog(context.Background(),
			"search index error while searching for an event", err.Error())
		if err != nil {
			logger.Error("can't send log to manager", "err", err)
		}

		return nil, err
	}
	var finalResult []event.Event

	for _, hit := range searchResult.Hits {

		hitJSON, err := json.Marshal(hit)
		if err != nil {
			_, err := h.grpc.AddLog(context.Background(),
				"error marshaling search result:", err.Error())
			if err != nil {
				logger.Error("can't send log to manager", "err", err)
			}
			continue
		}

		var newEvent event.Event
		if err := json.Unmarshal(hitJSON, &newEvent); err != nil {
			_, err := h.grpc.AddLog(context.Background(),
				"can't unmarshal search result to event:", err.Error())
			if err != nil {
				logger.Error("can't send log to manager", "err", err)
			}
			continue
		}

		if newEvent.Kind == types.KindGiftWrap {
			if newEvent.Tags.ContainsTag("p", pubkey) == false {
				continue // exclude others gift wrap events from final result
			}
		}

		finalResult = append(finalResult, newEvent)
	}

	return finalResult, nil
}
