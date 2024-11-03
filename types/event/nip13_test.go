package event_test

import (
	"testing"

	"github.com/dezh-tech/immortal/types/event"
	"github.com/stretchr/testify/assert"
)

func TestDifficulty(t *testing.T) {
	testCases := []struct {
		result int
		id     string
	}{
		{36, "000000000e9d97a1ab09fc381030b346cdd7a142ad57e6df0b46dc9bef6c7e2d"},
		{22, "0000024d38993bae75a61e82710842305fac9cda280f541476c31426c42ca81a"},
		{0, "f2775f4eeaa0aa45f66440490b45d6aede8d1ceb7ac443e6328763db5ce8d6e3"},
		{7, "010b807e82a1417588be0bcd7606b2aae4163a365afe1f6c97404b17fc56d30b"},
		{21, "000004f20d022b65dd961cbbdc157347dbd37ca375899fe38b40d174cccd8ee3"},
		{18, "00003db72b8385511ef2c1dd5fb3a43988269c9d7f51986f03c1aecd675dc506"},
	}

	for _, tc := range testCases {
		e := event.Event{ID: tc.id}
		assert.Equal(t, tc.result, e.Difficulty())
	}
}
