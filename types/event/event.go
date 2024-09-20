package event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/dezh-tech/immortal/types"
	"github.com/mailru/easyjson"
)

// Event represents an event structure defined on NIP-01.
type Event struct {
	ID        string     `json:"id"`
	PublicKey string     `json:"pubkey"`
	CreatedAt int64      `json:"created_at"`
	Kind      types.Kind `json:"kind"`
	Tags      types.Tags `json:"tags"`
	Content   string     `json:"content"`
	Signature string     `json:"sig"`
}

// Decode decodes a byte array into Event structure.
func Decode(b []byte) (*Event, error) {
	e := new(Event)

	if err := easyjson.Unmarshal(b, e); err != nil {
		return nil, types.DecodeError{
			Reason: err.Error(),
		}
	}

	return e, nil
}

// Encode encodes an Event to a byte array.
func (e *Event) Encode() ([]byte, error) {
	b, err := easyjson.Marshal(e)
	if err != nil {
		return nil, types.EncodeError{
			Reason: err.Error(),
		}
	}

	return b, nil
}

func (e *Event) Serialize() []byte {
	// the serialization process is just putting everything into a JSON array
	// so the order is kept. See NIP-01
	dst := make([]byte, 0)

	// the header portion is easy to serialize.
	// [0,"pubkey",created_at,kind,[
	dst = append(dst, []byte(
		fmt.Sprintf( //nolint
			"[0,\"%s\",%d,%d,",
			e.PublicKey,
			e.CreatedAt,
			e.Kind,
		))...)

	// tags.
	dst = types.MarshalTo(e.Tags, dst)
	dst = append(dst, ',')

	// content needs to be escaped in general as it is user generated.
	dst = types.EscapeString(dst, e.Content)
	dst = append(dst, ']')

	return dst
}

// GetID calculates ID of a given event.
func (e *Event) GetRawID() [32]byte {
	return sha256.Sum256(e.Serialize())
}

// IsValid function validats an event Signature and ID.
func (e *Event) IsValid(id [32]byte) bool {
	pk, err := hex.DecodeString(e.PublicKey)
	if err != nil {
		return false
	}

	pubkey, err := schnorr.ParsePubKey(pk)
	if err != nil {
		return false
	}

	s, err := hex.DecodeString(e.Signature)
	if err != nil {
		return false
	}

	sig, err := schnorr.ParseSignature(s)
	if err != nil {
		return false
	}

	// TODO::: replace with libsecp256k1 (C++ version).
	return sig.Verify(id[:], pubkey)
}

// String returns and encoded string representation of event e.
func (e *Event) String() string {
	ee, err := e.Encode()
	if err != nil {
		return ""
	}

	return string(ee)
}
