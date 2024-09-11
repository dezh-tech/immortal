package message_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name             string
	Message          []byte
	ExpectedEnvelope message.Message
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
		Name:    "REQ message",
		Message: []byte(`["REQ","million", {"kinds": [1]}, {"kinds": [30023 ], "#d": ["buteko",    "batuke"]}]`),
		ExpectedEnvelope: &message.Req{
			SubscriptionID: "million",
			Filters: filter.Filters{{Kinds: []types.Kind{1}}, {
				Kinds: []types.Kind{30023},
				Tags:  map[string]types.Tag{"d": []string{"buteko", "batuke"}},
			}},
		},
	},
	{
		Name:    "EVENT message",
		Message: []byte(`["EVENT",{"kind":1,"id":"d86745e397de3a3b8c8f15f7a02c6aa6a60213f59b14e1a7093ea286e020423e","pubkey":"79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798","created_at":1726055814,"tags":[],"content":"test","sig":"ec9f5702c2698ebfa1ce037c31568c7eb420d65d9684f33ba4d3266c82b771f58cf92c44f1e2a710ba71ce3c92bba1253aa7419d124b2a5eed74e9165e868d50"}]`),
		ExpectedEnvelope: &message.Event{
			Event: &event.Event{
				ID:        "d86745e397de3a3b8c8f15f7a02c6aa6a60213f59b14e1a7093ea286e020423e",
				PublicKey: "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
				CreatedAt: 1726055814,
				Kind:      types.KindTextNote,
				Tags:      []types.Tag{},
				Content:   "test",
				Signature: "ec9f5702c2698ebfa1ce037c31568c7eb420d65d9684f33ba4d3266c82b771f58cf92c44f1e2a710ba71ce3c92bba1253aa7419d124b2a5eed74e9165e868d50",
			},
		},
	},
}

func TestEnvelope(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			parsedEnvelope := message.ParseMessage(tc.Message)

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
