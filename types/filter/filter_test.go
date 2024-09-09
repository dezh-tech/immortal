package filter_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	RawFilter    string
	FilterObject *filter.Filter
	IsValidData  bool
	Note         string
}

var (
	testCases = []TestCase{
		{
			RawFilter: `{"ids": ["abc"],"#e":["zzz"],"limit":0,"#something":["nothing","bab"],"since":1644254609,"search":"test"}`,
			FilterObject: &filter.Filter{
				IDs: []string{"abc"},
				Tags: map[string]types.Tag{
					"e":         {"zzz"},
					"something": {"nothing", "bab"},
				},
				Since:  1644254609,
				Search: "test",
			},
			IsValidData: true,
			Note:        "this is a valid filter.",
		},
		{
			RawFilter: `{"kinds":[1],"authors":["a8171781fd9e90ede3ea44ddca5d3abf828fe8eedeb0f3abb0dd3e563562e1fc","1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59","ed4ca520e9929dfe9efdadf4011b53d30afd0678a09aa026927e60e7a45d9244"],"since":1677033299}`,
			FilterObject: &filter.Filter{
				Kinds: []types.Kind{types.KindTextNote},
				Authors: []string{
					"a8171781fd9e90ede3ea44ddca5d3abf828fe8eedeb0f3abb0dd3e563562e1fc",
					"1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59",
					"ed4ca520e9929dfe9efdadf4011b53d30afd0678a09aa026927e60e7a45d9244",
				},
				Tags:  map[string]types.Tag{},
				Since: 1677033299,
			},
			IsValidData: true,
			Note:        "this is a valid filter.",
		},
		{
			RawFilter: `{"kinds":[1,2,4],"until":12345678,"limit":0,"#fruit":["banana","mango"]}`,
			FilterObject: &filter.Filter{
				Kinds: []types.Kind{1, 2, 4},
				Until: 12345678,
				Limit: 0,
				Tags: map[string]types.Tag{
					"fruit": {"banana", "mango"},
				},
			},
			IsValidData: true,
			Note:        "this is a valid filter.",
		},
		{
			RawFilter: `{"kinds":[1,2,4],"until":12345678,"#fruit":["banana","mango"]}`,
			FilterObject: &filter.Filter{
				Kinds: []types.Kind{1, 2, 4},
				Until: 12345678,
				Tags: map[string]types.Tag{
					"fruit": {"banana", "mango"},
				},
			},
			IsValidData: true,
			Note:        "this is a valid filter.",
		},
		{
			RawFilter:    `{kinds":[1,2,4],"until":12345678#fruit":["banana","mango"]}`,
			FilterObject: &filter.Filter{},
			IsValidData:  false,
			Note:         "this is an invalid filter.",
		},
	}

	// TODO::: Add more test cases for matchs.
	testFilter = `{"kinds":[1],"authors":["1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59"],"since":1677033299}`
	testEvent  = `{"id":"5a127c9c931f392f6afc7fdb74e8be01c34035314735a6b97d2cf360d13cfb94","pubkey":"1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59","created_at":1677033299,"kind":1,"tags":[["t","japan"]],"content":"If you like my art,I'd appreciate a coin or two!!\nZap is welcome!! Thanks.\n\n\n#japan #bitcoin #art #bananaart\nhttps://void.cat/d/CgM1bzDgHUCtiNNwfX9ajY.webp","sig":"828497508487ca1e374f6b4f2bba7487bc09fccd5cc0d1baa82846a944f8c5766918abf5878a580f1e6615de91f5b57a32e34c42ee2747c983aaf47dbf2a0255"}`
)

func TestFilter(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		for i, tc := range testCases {
			f, err := filter.Decode([]byte(tc.RawFilter))
			if tc.IsValidData {
				assert.NoError(t, err, tc.Note)

				assert.Equal(t, tc.FilterObject.Authors, f.Authors)
				assert.Equal(t, tc.FilterObject.IDs, f.IDs)
				assert.Equal(t, tc.FilterObject.Kinds, f.Kinds)
				assert.Equal(t, tc.FilterObject.Limit, f.Limit)
				assert.Equal(t, tc.FilterObject.Search, f.Search)
				assert.Equal(t, tc.FilterObject.Since, f.Since)
				assert.Equal(t, tc.FilterObject.Tags, f.Tags, tc.RawFilter, "\n", i)
				assert.Equal(t, tc.FilterObject.Until, f.Until)

				continue
			}

			assert.Error(t, err, tc.Note)
		}
	})

	t.Run("Encode", func(t *testing.T) {
		for _, tc := range testCases {
			if tc.IsValidData {
				_, err := tc.FilterObject.Encode()
				assert.NoError(t, err, tc.Note)
			}
		}
	})

	t.Run("Match", func(t *testing.T) {
		e, err := event.Decode([]byte(testEvent))
		require.NoError(t, err)

		f, err := filter.Decode([]byte(testFilter))
		require.NoError(t, err)

		assert.True(t, f.Match(e))
	})
}
