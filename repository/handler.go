package repository

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dezh-tech/immortal/infrastructure/database"
	grpcclient "github.com/dezh-tech/immortal/infrastructure/grpc_client"
	"github.com/dezh-tech/immortal/infrastructure/meilisearch"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/filter"
)

type Handler struct {
	db     *database.Database
	meili  *meilisearch.Meili
	grpc   grpcclient.IClient
	config Config
}

func New(cfg Config, db *database.Database, meili *meilisearch.Meili, grpc grpcclient.IClient) *Handler {
	return &Handler{
		db:     db,
		meili:  meili,
		config: cfg,
		grpc:   grpc,
	}
}

func buildMeiliQuery(f *filter.Filter) string {
	var filters []string

	filters = append(filters, "pubkey IS NOT NULL")

	if len(f.IDs) > 0 {
		ids := strings.Join(f.IDs, "\", \"")
		filters = append(filters, fmt.Sprintf("id IN [%q]", ids))
	}

	if len(f.Kinds) > 0 {
		kindStrings := make([]string, len(f.Kinds))
		for i, kind := range f.Kinds {
			kindStrings[i] = fmt.Sprintf("%d", kind)
		}
		filters = append(filters, fmt.Sprintf("kind IN [%s]", strings.Join(kindStrings, ", ")))
	}

	if len(f.Authors) > 0 {
		quotedAuthors := make([]string, len(f.Authors))
		for i, author := range f.Authors {
			quotedAuthors[i] = fmt.Sprintf("%q", author)
		}
		filters = append(filters, fmt.Sprintf("pubkey IN [%s]", strings.Join(quotedAuthors, ", ")))
	}

	if len(f.Tags) > 0 {
		for tagKey, tagValues := range f.Tags {
			var tagConditions []string

			for _, tagValue := range tagValues {
				tagConditions = append(tagConditions,
					fmt.Sprintf("tags CONTAINS %q AND tags CONTAINS %q", tagKey, tagValue))
			}

			newFilter := fmt.Sprintf("(%s)", strings.Join(tagConditions, " OR "))
			filters = append(filters, newFilter)
		}
	}

	if f.Since > 0 {
		filters = append(filters, fmt.Sprintf("created_at >= %d", f.Since))
	}

	if f.Until > 0 {
		filters = append(filters, fmt.Sprintf("created_at <= %d", f.Until))
	}

	return strings.Join(filters, " AND ")
}

func filterToMongoQuery(f *filter.Filter, isMultiKindColl bool, k types.Kind, pubkey string) bson.D {
	query := make(bson.D, 0)

	if isMultiKindColl {
		query = append(query, bson.E{Key: "kind", Value: k})
	}

	query = append(query, bson.E{Key: "pubkey", Value: bson.M{
		"$exists": true,
	}})

	if len(f.IDs) > 0 {
		query = append(query, bson.E{Key: "id", Value: bson.M{"$in": f.IDs}})
	}

	if len(f.Authors) > 0 {
		query = append(query, bson.E{Key: "pubkey", Value: bson.M{"$in": f.Authors}})
	}

	if k == types.KindGiftWrap && pubkey != "" {
		f.Tags["p"] = []string{pubkey}
	}

	if len(f.Tags) > 0 {
		tagQueries := bson.A{}
		for tagKey, tagValues := range f.Tags {
			qtf := bson.M{
				"tags": bson.M{
					"$elemMatch": bson.M{
						"0": tagKey,
						"1": bson.M{"$in": tagValues},
					},
				},
			}
			tagQueries = append(tagQueries, qtf)
		}
		query = append(query, bson.E{Key: "$and", Value: tagQueries})
	}

	if f.Since > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$gte": f.Since}})
	}

	if f.Until > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$lte": f.Until}})
	}

	return query
}

func getCollectionName(k types.Kind) (string, bool) {
	collName, ok := types.KindToName[k]
	if ok {
		return collName, false
	}

	if k >= 9000 && k <= 9030 {
		return "groups", true
	}

	if k >= 1630 && k <= 1633 {
		return "status", true
	}

	if k >= 39000 && k <= 39009 {
		return "groups_metadata", true
	}

	if k >= 5000 && k <= 5999 || k >= 6000 && k <= 6999 || k == 7000 {
		return "dvm", true
	}

	return "unknown", true
}

func removeDuplicateKinds(intSlice []types.Kind) []types.Kind {
	allKeys := make(map[types.Kind]bool, len(intSlice))
	list := []types.Kind{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	return list
}
