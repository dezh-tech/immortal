package event_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	RawEvent    string
	EventObject *event.Event
	IsValidSig  bool
	IsValidData bool
	Note        string
}

var testCases = []TestCase{
	{
		RawEvent: `{"id":"a1d7ba3cdcc67a358186f85e5f2a02abd173877d484b76d1f1f22ee47d68293d","pubkey":"32e1827635450ebb3c5a7d12c1f8e7b2b514439ac10a67eef3d9fd9c5c68e245","created_at":1725890895,"kind":1,"tags":[],"content":"ReplyGuy never replies to me :( i feel left out","sig":"c2e6975905e41837343dc4b607dadf2895df457a0b8461b0f86d25506c4458c3fe83ed1f924715a0416412858fa5c51f3f3271361d729037f18d216b29618dda"}`,
		EventObject: &event.Event{
			ID:        "a1d7ba3cdcc67a358186f85e5f2a02abd173877d484b76d1f1f22ee47d68293d",
			PublicKey: "32e1827635450ebb3c5a7d12c1f8e7b2b514439ac10a67eef3d9fd9c5c68e245",
			CreatedAt: 1725890895,
			Kind:      types.KindTextNote,
			Tags:      []types.Tag{},
			Content:   "ReplyGuy never replies to me :( i feel left out",
			Signature: "c2e6975905e41837343dc4b607dadf2895df457a0b8461b0f86d25506c4458c3fe83ed1f924715a0416412858fa5c51f3f3271361d729037f18d216b29618dda",
		},
		IsValidSig:  true,
		IsValidData: true,
		Note:        "this is a valid event.",
	},
	{
		RawEvent: `{"content":"SUPER DOWN nostr:npub1h8nk2346qezka5cpm8jjh3yl5j88pf4ly2ptu7s6uu55wcfqy0wq36rpev","created_at":1725877943,"id":"a93df9f6746dfbd4de63196a36a0aa408dec8308fb55b3d5edcd22c953a4efb9","kind":1,"pubkey":"472be9f9264eea1254f2b2f7cd2da0c319dae4fe4cd649f0424e94234dcacf97","sig":"dd8e9478c52d086a793084dd092684344d3946b5f7f537573076530c07225870b732fa0e7214b0754e943e90340335756c8de5ab7a61c6f1375a4e5f340b6a26","tags":[["e","f6e8673a61ade88c087f45a6fa4f278e6e8b78dad2512a43b9e5a82e6df4ade4","","root"],["e","cec30d76c2599215b09f67f5d65f3b2786390bcfb51c97e8dab14c3d4e3c4e73","wss://relay.primal.net","reply"],["p","b9e76546ba06456ed301d9e52bc49fa48e70a6bf2282be7a1ae72947612023dc"],["p","b9e76546ba06456ed301d9e52bc49fa48e70a6bf2282be7a1ae72947612023dc"]]}`,
		EventObject: &event.Event{
			ID:        "a93df9f6746dfbd4de63196a36a0aa408dec8308fb55b3d5edcd22c953a4efb9",
			PublicKey: "472be9f9264eea1254f2b2f7cd2da0c319dae4fe4cd649f0424e94234dcacf97",
			CreatedAt: 1725877943,
			Kind:      types.KindTextNote,
			Tags: []types.Tag{
				{"e", "f6e8673a61ade88c087f45a6fa4f278e6e8b78dad2512a43b9e5a82e6df4ade4", "", "root"},
				{
					"e",
					"cec30d76c2599215b09f67f5d65f3b2786390bcfb51c97e8dab14c3d4e3c4e73",
					"wss://relay.primal.net",
					"reply",
				},
				{
					"p",
					"b9e76546ba06456ed301d9e52bc49fa48e70a6bf2282be7a1ae72947612023dc",
				},
				{
					"p",
					"b9e76546ba06456ed301d9e52bc49fa48e70a6bf2282be7a1ae72947612023dc",
				},
			},
			Content:   "SUPER DOWN nostr:npub1h8nk2346qezka5cpm8jjh3yl5j88pf4ly2ptu7s6uu55wcfqy0wq36rpev",
			Signature: "dd8e9478c52d086a793084dd092684344d3946b5f7f537573076530c07225870b732fa0e7214b0754e943e90340335756c8de5ab7a61c6f1375a4e5f340b6a26",
		},
		IsValidSig:  true,
		IsValidData: true,
		Note:        "this is a valid event.",
	},
	{
		RawEvent: `{"content":"that’s a link to another website.","created_at":1725832414,"id":"594915a98c7f65b65a642e076463f5ac1319ae55d116c401528094e56023abf8","kind":1,"pubkey":"472be9f9264eea1254f2b2f7cd2da0c319dae4fe4cd649f0424e94234dcacf97","sig":"fc6d27f4cf775d7190a597fe67f7b0a4341ad044ae29d7d2b8daab060cf4e600a89521cbf0c370e00c8252b0ef4553294ebb56794cabe47ec7c16a800d857ca6","tags":[["e","b23ea78bd672a85faa84cdf4206231c217ae9034b9d82be78528279ffc87bbb9","","root"],["e","6ec00066a620631fb396d1756bdf56b397ea5626d88e71a8ec533e6ce99d595f","wss://a.nos.lol","reply"],["p","63fe6318dc58583cfe16810f86dd09e18bfd76aabc24a0081ce2856f330504ed"]]}`,
		EventObject: &event.Event{
			ID:        "594915a98c7f65b65a642e076463f5ac1319ae55d116c401528094e56023abf8",
			PublicKey: "472be9f9264eea1254f2b2f7cd2da0c319dae4fe4cd649f0424e94234dcacf97",
			CreatedAt: 1725832414,
			Kind:      types.KindTextNote,
			Tags: []types.Tag{
				{
					"e",
					"b23ea78bd672a85faa84cdf4206231c217ae9034b9d82be78528279ffc87bbb9",
					"",
					"root",
				},
				{
					"e",
					"6ec00066a620631fb396d1756bdf56b397ea5626d88e71a8ec533e6ce99d595f",
					"wss://a.nos.lol",
					"reply",
				},
				{
					"p",
					"63fe6318dc58583cfe16810f86dd09e18bfd76aabc24a0081ce2856f330504ed",
				},
			},
			Content:   "that’s a link to another website.",
			Signature: "fc6d27f4cf775d7190a597fe67f7b0a4341ad044ae29d7d2b8daab060cf4e600a89521cbf0c370e00c8252b0ef4553294ebb56794cabe47ec7c16a800d857ca6",
		},
		IsValidSig:  false,
		IsValidData: true,
		Note:        "sig field is invalid.",
	},
	{
		RawEvent:    `"content:"SUPER DO "np7e8dab14c3d4e3c4e73","wss://relay.primal.net","reply"],["p","b9e76546ba06456ed301d9e52bc49fa48e70a6bf2282be7a1ae72947612023dc"],["p","b9e76546ba06456ed301d9e52bc49fa48e70a6bf2282be7a1ae72947612023dc"]]}`,
		EventObject: nil,
		IsValidSig:  true,
		IsValidData: false,
		Note:        "data and encoding are invalid.",
	},
}

func TestEvent(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		for _, tc := range testCases {
			e, err := event.Decode([]byte(tc.RawEvent))
			if tc.IsValidData {
				assert.NoError(t, err, tc.Note)

				assert.Equal(t, tc.EventObject.ID, e.ID)
				assert.Equal(t, tc.EventObject.CreatedAt, e.CreatedAt)
				assert.Equal(t, tc.EventObject.Content, e.Content)
				assert.Equal(t, tc.EventObject.Kind, e.Kind)
				assert.Equal(t, tc.EventObject.PublicKey, e.PublicKey)
				assert.Equal(t, tc.EventObject.Signature, e.Signature)

				continue
			}

			assert.Error(t, err, tc.Note)
		}
	})

	t.Run("Encode", func(t *testing.T) {
		for _, tc := range testCases {
			if tc.IsValidData {
				e, err := tc.EventObject.Encode()
				assert.NoError(t, err, tc.Note)
				assert.Equal(t, len([]byte(tc.RawEvent)), len(e))

				// assert.Equal(t, tc.RawEvent, string(e)) //TODO:: is that correct?
			}
		}
	})

	t.Run("CheckSig", func(t *testing.T) {
		for _, tc := range testCases {
			if tc.IsValidData {
				isValid := tc.EventObject.IsValid(tc.EventObject.GetRawID())
				if tc.IsValidSig {
					assert.True(t, isValid, tc.Note)

					continue
				}

				assert.False(t, isValid, tc.Note)
			}
		}
	})
}
