package envelope_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/envelope"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name             string
	Message          []byte
	ExpectedEnvelope envelope.Envelope
}

var testCases = []testCase{
	{
		Name:             "nil",
		Message:          nil,
		ExpectedEnvelope: nil,
	},
	{
		Name:             "invalid string",
		Message:          []byte("invalid input"),
		ExpectedEnvelope: nil,
	},
	{
		Name:             "invalid string with a comma",
		Message:          []byte("invalid, input"),
		ExpectedEnvelope: nil,
	},
	{
		Name:    "REQ envelope",
		Message: []byte(`["REQ","million", {"kinds": [1]}, {"kinds": [30023 ], "#d": ["buteko",    "batuke"]}]`),
		ExpectedEnvelope: &envelope.ReqEnvelope{
			SubscriptionID: "million",
			Filters: filter.Filters{{Kinds: []types.Kind{1}}, {
				Kinds: []types.Kind{30023},
				Tags:  map[string]types.Tag{"d": []string{"buteko", "batuke"}},
			}},
		},
	},
}

func TestEnvelope(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			parsedEnvelope := envelope.ParseMessage(tc.Message)

			if tc.ExpectedEnvelope == nil && parsedEnvelope == nil {
				return
			}

			if tc.ExpectedEnvelope == nil {
				assert.NotNil(t, parsedEnvelope)
			}

			assert.Equal(t, tc.ExpectedEnvelope.String(), parsedEnvelope.String())
		})
	}
}
