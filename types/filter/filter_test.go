package filter_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types/filter"
	"github.com/stretchr/testify/assert"
)

// TODO::: Use table test.
// TODO::: Add error cases.

var (
	rawValidFilters = []string{
		`{"ids": ["abc"],"#e":["zzz"],"#something":["nothing","bab"],"since":1644254609,"search":"test"}`,
		`{"ids": ["abc"],"#e":["zzz"],"limit":0,"#something":["nothing","bab"],"since":1644254609,"search":"test"}`,
		`{"kinds":[1],"authors":["a8171781fd9e90ede3ea44ddca5d3abf828fe8eedeb0f3abb0dd3e563562e1fc","1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59","ed4ca520e9929dfe9efdadf4011b53d30afd0678a09aa026927e60e7a45d9244"],"since":1677033299}`,
		`{"kinds":[1,2,4],"until":12345678,"limit":0,"#fruit":["banana","mango"]}`,
		`{"kinds":[1,2,4],"until":12345678,"#fruit":["banana","mango"]}`,
	}

	EncodedFilter []byte
	DecodedFilter *filter.Filter
)

func TestDencode(t *testing.T) {
	for _, f := range rawValidFilters {
		_, err := filter.Decode([]byte(f))
		assert.NoError(t, err)
	}
}

func BenchmarkDecode(b *testing.B) {
	var decodedFilter *filter.Filter
	for i := 0; i < b.N; i++ {
		for _, f := range rawValidFilters {
			decodedFilter, _ = filter.Decode([]byte(f))
		}
	}
	DecodedFilter = decodedFilter
}

func TestEncode(t *testing.T) {
	filters := make([]*filter.Filter, len(rawValidFilters))

	for _, f := range rawValidFilters {
		decodedFilter, err := filter.Decode([]byte(f))
		assert.NoError(t, err)

		filters = append(filters, decodedFilter)
	}

	for _, f := range filters {
		_, err := f.Encode()
		assert.NoError(t, err)
	}
}

func BenchmarkEncode(b *testing.B) {
	filters := make([]*filter.Filter, len(rawValidFilters))

	for _, f := range rawValidFilters {
		decodedFilter, _ := filter.Decode([]byte(f))
		filters = append(filters, decodedFilter)
	}

	b.ResetTimer()

	var encodedFilter []byte
	for i := 0; i < b.N; i++ {
		for _, f := range filters {
			encodedFilter, _ = f.Encode()
		}
	}
	EncodedFilter = encodedFilter
}
