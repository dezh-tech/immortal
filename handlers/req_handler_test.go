package handlers

import (
	"testing"
	"time"

	"github.com/dezh-tech/immortal/types"
)

func TestBuildDynamicQueries(t *testing.T) {
	types.KindToTableMap = map[types.Kind]string{
		types.KindTextNote: "text_notes",
		types.KindReaction: "reactions",
	}

	tests := []struct {
		name     string
		input    map[types.Kind][]queryFilter
		expected string
		args     []interface{}
		wantErr  bool
	}{
		{
			name: "Simple Case - One Filter with IDs",
			input: map[types.Kind][]queryFilter{
				types.KindTextNote: {
					{
						IDs:   []string{"id1", "id2"},
						Limit: 10,
					},
				},
			},
			expected: "SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE id = ANY($1) LIMIT 10) AS queryfilter1 ORDER BY event_created_at ASC, id ASC",
			args:     []interface{}{[]string{"id1", "id2"}},
			wantErr:  false,
		},
		{
			name: "Case with Authors, Since, and Until",
			input: map[types.Kind][]queryFilter{
				types.KindTextNote: {
					{
						Authors: []string{"author1", "author2"},
						Since:   time.Now().Add(-time.Hour * 24).Unix(),
						Until:   time.Now().Unix(),
					},
				},
			},
			expected: "SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE users_metadatapub_key = ANY($1) AND created_at >= $2 AND created_at <= $3 LIMIT 20) AS queryfilter1 ORDER BY event_created_at ASC, id ASC",
			args:     []interface{}{[]string{"author1", "author2"}, time.Unix(time.Now().Add(-time.Hour*24).Unix(), 0), time.Unix(time.Now().Unix(), 0)},
			wantErr:  false,
		},
		{
			name: "No Filters Provided",
			input: map[types.Kind][]queryFilter{
				types.KindTextNote: {
					{
						IDs:   []string{},
						Limit: 10,
					},
				},
			},
			expected: "",
			args:     nil,
			wantErr:  true,
		},
		{
			name: "Filter with Tags",
			input: map[types.Kind][]queryFilter{
				types.KindTextNote: {
					{
						Tags: map[string]types.Tag{
							"e": {"event_tag_value1", "event_tag_value2"},
							"p": {"person_tag_value1", "person_tag_value2"},
						},
						Limit: 15,
					},
				},
			},
			expected: "SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE e_tags = ANY($1) AND p_tags = ANY($2) LIMIT 15) AS queryfilter1 ORDER BY event_created_at ASC, id ASC",
			args:     []interface{}{[]string{"event_tag_value1", "event_tag_value2"}, []string{"person_tag_value1", "person_tag_value2"}},
			wantErr:  false,
		},
		{
			name: "Error Case - Invalid Kind without Table Mapping",
			input: map[types.Kind][]queryFilter{
				types.Kind(999): { // 999 is an invalid kind
					{
						IDs:   []string{"id1"},
						Limit: 5,
					},
				},
			},
			expected: "",
			args:     []interface{}{},
			wantErr:  true,
		},
		{
			name: "Multiple Filters - Generate multiple queryfilters",
			input: map[types.Kind][]queryFilter{
				types.KindTextNote: {
					{
						IDs:   []string{"id1"},
						Limit: 10,
					},
					{
						Authors: []string{"author1"},
						Limit:   5,
					},
				},
			},
			expected: "SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE id = ANY($1) LIMIT 10) AS queryfilter1 UNION " +
				"SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE users_metadatapub_key = ANY($2) LIMIT 5) AS queryfilter2 ORDER BY event_created_at ASC, id ASC",
			args:    []interface{}{[]string{"id1"}, []string{"author1"}},
			wantErr: false,
		},
		{
			name: "Two Different Tables, Multiple Subqueries for Each",
			input: map[types.Kind][]queryFilter{
				types.KindReaction: {
					{
						IDs:   []string{"reaction_id1", "reaction_id2"},
						Limit: 15,
					},
					{
						Authors: []string{"reaction_author1"},
						Limit:   10,
					},
					{
						Tags: map[string]types.Tag{
							"e": {"reaction_tag_value1", "reaction_tag_value2"},
						},
						Limit: 20,
					},
				},
				types.KindTextNote: {
					{
						IDs:   []string{"text_id1"},
						Limit: 10,
					},
					{
						Authors: []string{"author1"},
						Limit:   5,
					},
				},
			},
			expected: "SELECT * FROM (SELECT id, event_created_at, event FROM reactions WHERE id = ANY($1) LIMIT 15) AS queryfilter1 UNION SELECT * FROM (SELECT id, event_created_at, event FROM reactions WHERE users_metadatapub_key = ANY($2) LIMIT 10) AS queryfilter2 UNION SELECT * FROM (SELECT id, event_created_at, event FROM reactions WHERE e_tags = ANY($3) LIMIT 20) AS queryfilter3 UNION SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE id = ANY($4) LIMIT 10) AS queryfilter4 UNION SELECT * FROM (SELECT id, event_created_at, event FROM text_notes WHERE users_metadatapub_key = ANY($5) LIMIT 5) AS queryfilter5 ORDER BY event_created_at ASC, id ASC",
			args: []interface{}{
				[]string{"text_id1"},
				[]string{"author1"},
				[]string{"reaction_id1", "reaction_id2"},
				[]string{"reaction_author1"},
				[]string{"reaction_tag_value1", "reaction_tag_value2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := BuildDynamicQueries(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildDynamicQueries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if query != tt.expected {
				t.Errorf("BuildDynamicQueries() query = %v, expected %v", query, tt.expected)
			}
			// TODO :::Compare args element by element
		})
	}
}