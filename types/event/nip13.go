package event

import (
	"encoding/hex"
	"math/bits"
)

// Difficulty returns the leading zeros of event id in base-2.
func (e *Event) Difficulty() int {
	var zeros int
	var b [1]byte

	for i := 0; i < 64; i += 2 {
		if e.ID[i:i+2] == "00" {
			zeros += 8

			continue
		}

		if _, err := hex.Decode(b[:], []byte{e.ID[i], e.ID[i+1]}); err != nil {
			return -1
		}

		zeros += bits.LeadingZeros8(b[0])

		break
	}

	return zeros
}
