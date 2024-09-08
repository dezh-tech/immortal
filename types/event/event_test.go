package event_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/stretchr/testify/assert"
)

// TODO::: Use table test.
// TODO::: Add error cases.

var (
	validRawEvents = []string{
		`{"kind":1,"id":"dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962","pubkey":"3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d","created_at":1644271588,"tags":[],"content":"now that https://blueskyweb.org/blog/2-7-2022-overview was announced we can stop working on nostr?","sig":"230e9d8f0ddaf7eb70b5f7741ccfa37e87a455c9a469282e3464e2052d3192cd63a167e196e381ef9d7e69e9ea43af2443b839974dc85d8aaab9efe1d9296524"}`,
		`{"kind":1,"id":"dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962","pubkey":"3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d","created_at":1644271588,"tags":[],"content":"now that https://blueskyweb.org/blog/2-7-2022-overview was announced we can stop working on nostr?","sig":"230e9d8f0ddaf7eb70b5f7741ccfa37e87a455c9a469282e3464e2052d3192cd63a167e196e381ef9d7e69e9ea43af2443b839974dc85d8aaab9efe1d9296524","extrakey":55}`,
		`{"kind":1,"id":"dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962","pubkey":"3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d","created_at":1644271588,"tags":[],"content":"now that https://blueskyweb.org/blog/2-7-2022-overview was announced we can stop working on nostr?","sig":"230e9d8f0ddaf7eb70b5f7741ccfa37e87a455c9a469282e3464e2052d3192cd63a167e196e381ef9d7e69e9ea43af2443b839974dc85d8aaab9efe1d9296524","extrakey":"aaa"}`,
		`{"kind":3,"id":"9e662bdd7d8abc40b5b15ee1ff5e9320efc87e9274d8d440c58e6eed2dddfbe2","pubkey":"373ebe3d45ec91977296a178d9f19f326c70631d2a1b0bbba5c5ecc2eb53b9e7","created_at":1644844224,"tags":[["p","3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d"],["p","75fc5ac2487363293bd27fb0d14fb966477d0f1dbc6361d37806a6a740eda91e"],["p","46d0dfd3a724a302ca9175163bdf788f3606b3fd1bb12d5fe055d1e418cb60ea"]],"content":"{\"wss://nostr-pub.wellorder.net\":{\"read\":true,\"write\":true},\"wss://nostr.bitcoiner.social\":{\"read\":false,\"write\":true},\"wss://expensive-relay.fiatjaf.com\":{\"read\":true,\"write\":true},\"wss://relayer.fiatjaf.com\":{\"read\":true,\"write\":true},\"wss://relay.bitid.nz\":{\"read\":true,\"write\":true},\"wss://nostr.rocks\":{\"read\":true,\"write\":true}}","sig":"811355d3484d375df47581cb5d66bed05002c2978894098304f20b595e571b7e01b2efd906c5650080ffe49cf1c62b36715698e9d88b9e8be43029a2f3fa66be"}`,
		`{"id":"6ea18dd9156305d7716348a459683642b0a35693c301d99665dc0bd4c58872a2","pubkey":"bd4ae3e67e29964d494172261dc45395c89f6bd2e774642e366127171dfb81f5","content":"OK, this is a test case for Immortal.","kind":1,"created_at":1725802966,"tags":[],"sig":"4ec9407243f41ca0b1b44e3b61ca2e43d8a20ed088357b075ee123a720ecd9b1734526d79f3c97cd77828fe1176e37104cce1270eb739499fcb202fced766e72","relays":[]}`,
		`{"id":"6ea18dd9156305d7716348a459683642b0a35693c301d99665dc0bd4c58872a2","pubkey":"bd4ae3e67e29964d494172261dc45395c89f6bd2e774642e366127171dfb81f5","content":"OK, this is a test case for Immortal.","kind":1,"created_at":1725802966,"tags":[],"sig":"4ec9407243f41ca0b1b44e3b61ca2e43d8a20ed088357b075ee123a720ecd9b1734526d79f3c97cd77828fe1176e37104cce1270eb739499fcb202fced766e72","relays":[]}`,
		`{"content":"GOOD MORNING.\n\nLIVE FREE.\n\nhttps://cdn.satellite.earth/fbd7f2d73469c95ca7ef0f6f66cf9456c4dcce5cb46d69dcc1d3243fe817faf3.mp4","created_at":1725802688,"id":"b8b2d7f724e3e774226ba84a621155a3656b58baf08c12c56f5452fe71b4fec9","kind":1,"pubkey":"04c915daefee38317fa734444acee390a8269fe5810b2241e5e6dd343dfbecc9","sig":"c26537743dcd6a9adbec078d1394953ce55ca8c5af1b1d8a38ce716f946c5c6b59a3f9129b4df41ee537fcd332dbfddecd1126e30f65e1feff6dd426032a3e60","tags":[]}`,
		`{"content":"{\"id\":\"1807a9e2d51e8e04fa6257fb1c5746df57c83ac6127b4c6462f9d986d5c98736\",\"kind\":1,\"pubkey\":\"c48e29f04b482cc01ca1f9ef8c86ef8318c059e0e9353235162f080f26e14c11\",\"content\":\" https:\\/\\/i.nostr.build\\/llzsmQs5gOtF1r8d.jpg \",\"tags\":[[\"imeta\",\"url https:\\/\\/i.nostr.build\\/llzsmQs5gOtF1r8d.jpg\",\"blurhash enQlLDaJysoy-;nMWBozs.kWKQtRs8ayenx^oyW;V[enNxV@i^WBfl\",\"dim 1259x1259\"],[\"r\",\"https:\\/\\/i.nostr.build\\/llzsmQs5gOtF1r8d.jpg\"]],\"sig\":\"5612c5a2ee8224e6d4b386698b1a9ae137cb9fa1c91ba7f15c4a14ac7896550f82caeffab7376e48e6fb81acf17597959c7658cabcf30c7484aa84cd740f2fee\",\"created_at\":1725639394}","created_at":1725802348,"id":"d2c8db3990efc74a56c4e9602bdcd5763e586b24b0708c14a171923f5e36e184","kind":6,"pubkey":"c48e29f04b482cc01ca1f9ef8c86ef8318c059e0e9353235162f080f26e14c11","sig":"5d3fd5cff3628905178c829416f793304688596e0304b95dd94976f55885adf2afe13d6120d932867c40e240429cdd2c4c1930b41e2beea852ce4626f51fc35a","tags":[["e","1807a9e2d51e8e04fa6257fb1c5746df57c83ac6127b4c6462f9d986d5c98736","","root"],["p","c48e29f04b482cc01ca1f9ef8c86ef8318c059e0e9353235162f080f26e14c11"]]}`,
		`{"content":"Just add ReplyGuy now in Amethyst. \nnostr:nevent1qqsyjv4z7ns6frwnfrcn0lqk227chnrcnat476yaaadg8ev8jgn6p4gpz4mhxue69uhhyetvv9ujuerpd46hxtnfduhsygqpmhhz3xc696ggwn9rg2985s28vjnv45dtl25ctsspu74d59kn3spsgqqqqqqsa49ewc","created_at":1725799931,"id":"66e9072f5b2e7c4c6c6d8b9de2a81f78983549d2f56b7dc11737bba3a6a71408","kind":1,"pubkey":"01ddee289b1a2e90874ca3428a7a414764a6cad1abfaa985c201e7aada16d38c","sig":"2b17c30db9f7e246281d1dabcea12b40d19b6a6358d040a9b18a7b0fda72cca67adb260a7fe9b125a1ade5672cd922aa5a8d8d32fa4f18dc0efc80ae3f2deebe","tags":[["e","4932a2f4e1a48dd348f137fc1652bd8bcc789f575f689def5a83e5879227a0d5","","mention"],["p","01ddee289b1a2e90874ca3428a7a414764a6cad1abfaa985c201e7aada16d38c","","mention"],["q","4932a2f4e1a48dd348f137fc1652bd8bcc789f575f689def5a83e5879227a0d5"]]}`,
		`{"content":"Kid: I want to be your age\nMe: no you don’t","created_at":1725799866,"id":"d4516dd8eda1c6235f6ca919b3ad6bdc7fe0a97d9bc7388b0444d9362f166522","kind":1,"pubkey":"1bc70a0148b3f316da33fe3c89f23e3e71ac4ff998027ec712b905cd24f6a411","sig":"17ec79e0e82e4ac3b3a32cfeac745774da90983d293435df9fa3225f31db87d4c855b2aee9d1c51a9bd64308262c1c7bf4274af20798b14b085d4dfc8f21ef9c","tags":[]}`,
		`{"content":"Wen nostr hobby apps?","created_at":1725798384,"id":"43795c6b71168e6973223751d92f0904785a272eb40ae505243ae211cebddfa1","kind":1,"pubkey":"1bc70a0148b3f316da33fe3c89f23e3e71ac4ff998027ec712b905cd24f6a411","sig":"6b5309d3580b1f128279e66ea822cb153cfb7da14b1d130e86cb4855fb31b48d3f9ae6f5bd57732f96cf5914542c5f90d5fddb67e509fbd0efe534b41e568fdc","tags":[]}`,
		`{"content":"GM from cicada\n\nhttps://video.nostr.build/66373d8edea1fadaacd10f4c5590729fb6e80f3f2279254d4c6b16cdbd80797f.mp4","created_at":1725798022,"id":"75b0fbee92009f70f57770aaf6ca993619af032b38d4b9abfa9bfe8d79c0f933","kind":1,"pubkey":"0461fcbecc4c3374439932d6b8f11269ccdb7cc973ad7a50ae362db135a474dd","sig":"cf8bdc8e24e46c41aede59ca5a8766ead2c347b8a064b26e2ec30895a5f38ed6194c2bb8be11f6d17d5757de165364075c166ea1203777d299329f91ab245728","tags":[["imeta","url https://video.nostr.build/66373d8edea1fadaacd10f4c5590729fb6e80f3f2279254d4c6b16cdbd80797f.mp4","m video/mp4","x 6c0fe5d3347139a488f5f7ffcf179b5a61b28c65057eac951241d46c299db34d","ox 66373d8edea1fadaacd10f4c5590729fb6e80f3f2279254d4c6b16cdbd80797f","size 664508"]]}`,
	}

	DecodedEvent *event.Event
	EncodedEvent []byte

	events = []event.Event{
		{
			ID:        "dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962",
			PublicKey: "3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d",
			CreatedAt: 1644271588,
			Kind:      types.KindTextNote,
			Tags:      []types.Tag{},
			Content:   "now that https://blueskyweb.org/blog/2-7-2022-overview was announced we can stop working on nostr?",
			Signature: "230e9d8f0ddaf7eb70b5f7741ccfa37e87a455c9a469282e3464e2052d3192cd63a167e196e381ef9d7e69e9ea43af2443b839974dc85d8aaab9efe1d9296524",
		},
	}

	EventValidation bool
)

func TestDecode(t *testing.T) {
	for _, e := range validRawEvents {
		_, err := event.Decode([]byte(e))
		assert.NoError(t, err, "valid event must be decoded with no erros JSON")
	}
}

func BenchmarkDecode(b *testing.B) {
	var decodedEvent *event.Event
	for i := 0; i < b.N; i++ {
		for _, e := range validRawEvents {
			decodedEvent, _ = event.Decode([]byte(e))
		}
	}
	DecodedEvent = decodedEvent
}

func TestEncode(t *testing.T) {
	events := make([]*event.Event, len(validRawEvents))
	for _, e := range validRawEvents {
		decodedEvent, err := event.Decode([]byte(e))
		assert.NoError(t, err)

		events = append(events, decodedEvent)
	}

	for _, e := range events {
		_, err := e.Encode()
		assert.NoError(t, err)
	}
}

func BenchmarkEncode(b *testing.B) {
	events := make([]*event.Event, len(validRawEvents))
	for _, e := range validRawEvents {
		decodedEvent, _ := event.Decode([]byte(e))
		events = append(events, decodedEvent)
	}

	b.ResetTimer()

	var encodedEvent []byte
	for i := 0; i < b.N; i++ {
		for _, e := range events {
			encodedEvent, _ = e.Encode()
		}
	}
	EncodedEvent = encodedEvent
}

// TODO::: add more test cases + benchmark.
func TestMatch(t *testing.T) {
	f, err := filter.Decode([]byte(`{"kinds":[1, 2, 4],"authors":["3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d","1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59","ed4ca520e9929dfe9efdadf4011b53d30afd0678a09aa026927e60e7a45d9244"],"since":1677033299}`))
	assert.NoError(t, err)

	e, err := event.Decode([]byte(`{"kind":1,"id":"dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962","pubkey":"3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d","created_at":1644271588,"tags":[],"content":"now that https://blueskyweb.org/blog/2-7-2022-overview was announced we can stop working on nostr?","sig":"230e9d8f0ddaf7eb70b5f7741ccfa37e87a455c9a469282e3464e2052d3192cd63a167e196e381ef9d7e69e9ea43af2443b839974dc85d8aaab9efe1d9296524"}`))
	assert.NoError(t, err)

	assert.True(t, e.Match(*f))
}

func TestValidate(t *testing.T) {
	for _, e := range events {
		valid, err := e.IsValid()

		assert.NoError(t, err)
		assert.True(t, valid)
	}
}

func BenchmarkValidate(b *testing.B) {
	var eventValidation bool
	for i := 0; i < b.N; i++ {
		for _, e := range events {
			eventValidation, _ = e.IsValid()
		}
	}
	EventValidation = eventValidation
}
