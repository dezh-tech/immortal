package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dezh-tech/immortal/database"
	dbmodels "github.com/dezh-tech/immortal/database/models"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type ReqHandler struct {
	DB *database.Database
}

type queryFilter struct {
	Tags map[string]types.Tag

	Authors []string
	IDs     []string

	Since int64
	Until int64
	Limit uint16
}

type res struct{
	id string `boil:"id"`
	event_created_at time.Time `boil:"event_created_at"`
	event string `boil:"event"`
}

func NewReqHandler(db *database.Database) *ReqHandler {
	return &ReqHandler{
		DB: db,
	}
}

func (rh *ReqHandler) Handle(fs filter.Filters) {
	queryKinds := make(map[types.Kind][]queryFilter)

	for _, f := range fs {
		qf := queryFilter{
			Tags:    f.Tags,
			Authors: f.Authors,
			IDs:     f.IDs,
			Since:   f.Since,
			Until:   f.Until,
			Limit:   f.Limit,
		}

		uniqueKinds := removeDuplicateKind(f.Kinds)
		for _, k := range uniqueKinds {
			queryKinds[k] = append(queryKinds[k], qf)
		}
	}

	q, args, err := BuildDynamicQueries(queryKinds)
	if err != nil{
		//
	}

	var r res

	err = queries.RawG(q,args...).BindG(context.Background(), &r)
	if err != nil{
		//
	}

	

}

func removeDuplicateKind(intSlice []types.Kind) []types.Kind {
	allKeys := make(map[types.Kind]bool)
	list := []types.Kind{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func BuildDynamicQueries(a map[types.Kind][]queryFilter) (string, []interface{}, error) {
	var subQueries []string
	var args []interface{}
	i := 1 // For parameterized query
	x := 1 // For subQuery name

	for kind, qfs := range a {
		tableName, ok := types.KindToTableMap[kind]
		if !ok {
			continue
		}

		for _, qf := range qfs {

			var limit uint16
			if qf.Limit > 0 {
				limit = qf.Limit
			} else {
				limit = 20
			}

			var conditions []string //* id = ANY(), users_metadatapub_key = ANY() -> join with `AND`

			if len(qf.IDs) > 0 {
				conditions = append(conditions, fmt.Sprintf("id = ANY($%d)", i))
				args = append(args, qf.IDs)
				i++
			}

			if len(qf.Authors) > 0 {
				conditions = append(conditions, fmt.Sprintf("users_metadatapub_key = ANY($%d)", i))
				args = append(args, qf.Authors)
				i++
			}

			// TAGs
			if len(qf.Tags) > 0 {
				for tagKey, tagValue := range qf.Tags {
					switch tagKey {
					case "e":
						conditions = append(conditions, fmt.Sprintf("e_tags = ANY($%d)", i))
						args = append(args, tagValue)
						i++
					case "p":
						conditions = append(conditions, fmt.Sprintf("p_tags = ANY($%d)", i))
						args = append(args, tagValue)
						i++
					case "a":
						conditions = append(conditions, fmt.Sprintf("a_tags = ANY($%d)", i))
						args = append(args, tagValue)
						i++
					default:

					}
				}
			}

			if qf.Since > 0 {
				conditions = append(conditions, fmt.Sprintf("created_at >= $%d", i))
				args = append(args, time.Unix(qf.Since, 0))
				i++
			}

			if qf.Until > 0 {
				conditions = append(conditions, fmt.Sprintf("created_at <= $%d", i))
				args = append(args, time.Unix(qf.Until, 0))
				i++
			}

			if len(conditions) < 1 {
				continue
			}

			subQuery := fmt.Sprintf("SELECT * FROM (SELECT id, event_created_at, event FROM %s WHERE %s LIMIT %d) AS queryfilter%d", tableName, strings.Join(conditions, " AND "), limit, x)
			subQueries = append(subQueries, subQuery)
			x++
		}
	}

	if len(subQueries) == 0 {
		return "", nil, fmt.Errorf("no valid subqueries generated")
	}

	return fmt.Sprintf("%s ORDER BY event_created_at ASC, id ASC", strings.Join(subQueries, " UNION ")), args, nil
}
