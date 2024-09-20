package filter_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/stretchr/testify/assert"
)

type (
	encodingTestCase struct {
		RawFilter    string
		FilterObject *filter.Filter
		IsValidData  bool
		Note         string
	}

	matchingTestCase struct {
		Filter  *filter.Filter
		Event   *event.Event
		IsMatch bool
		Note    string
	}
)

var (
	encodingTestCases = []encodingTestCase{
		{
			RawFilter: `{"ids": ["abc"],"#e":["zzz"],"limit":0,"#something":["nothing","bab"],"since":1644254609,"search":"test"}`,
			FilterObject: &filter.Filter{
				IDs: []string{"abc"},
				Tags: map[string][]string{
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
				Tags:  map[string][]string{},
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
				Tags: map[string][]string{
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
				Tags: map[string][]string{
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

	matchingTestCases = []matchingTestCase{
		{
			Filter: &filter.Filter{
				IDs: []string{
					"97ff1da488407f9b62b139ca750b4d486132bc18105acf6b771acb36ae107900",
					"e641244e30894d6c0530b628232c1ca28954d1bb1bbf1c2b94b395ec8e672acc",
				},
				Kinds: []types.Kind{types.KindTextNote},
				Tags: map[string][]string{
					"#e": {"97ff1da488407f9b62b139ca750b4d486132bc18105acf6b771acb36ae107900"},
				},
			},
			Event: &event.Event{
				ID:        "e641244e30894d6c0530b628232c1ca28954d1bb1bbf1c2b94b395ec8e672acc",
				PublicKey: "ba0901be8b694476afe9cf55d861d437eab847431f2035a57f23382cc1f1ea34",
				CreatedAt: 1726804325,
				Kind:      types.KindTextNote,
				Tags: []types.Tag{
					{"e", "97ff1da488407f9b62b139ca750b4d486132bc18105acf6b771acb36ae107900"},
					{"e", "7c77fa4057f78ca892b4d5533dd46a3ec5d938631448a5bd65f345038bfe8f5d"},
					{
						"p", "ba0901be8b694476afe9cf55d861d437eab847431f2035a57f23382cc1f1ea34",
						"wss://relay.jellyfish.land",
					},
				},
				Content:   "test_content",
				Signature: "test_signature",
			},
			IsMatch: true,
			Note:    "this filter and event are match and valid!",
		},
	}
)

func TestFilter(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		for i, tc := range encodingTestCases {
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
		for _, tc := range encodingTestCases {
			if tc.IsValidData {
				_, err := tc.FilterObject.Encode()
				assert.NoError(t, err, tc.Note)
			}
		}
	})

	t.Run("Match", func(t *testing.T) {
		for _, tc := range matchingTestCases {
			if tc.IsMatch {
				assert.True(t, tc.Filter.Match(tc.Event),
					"expected event %s to be match with filter %s, note: %s",
					tc.Event.String(), tc.Filter.String(), tc.Note)
				continue
			}
			assert.False(t, tc.Filter.Match(tc.Event),
				"expected event %s not to be match with filter %s, note: %s",
				tc.Event.String(), tc.Filter.String(), tc.Note)
		}
	})
}
