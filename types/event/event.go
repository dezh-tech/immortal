package event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/dezh-tech/immortal/types"
	"github.com/mailru/easyjson"
)

const hextable = "0123456789abcdef"

// Event represents an event structure defined on NIP-01.
type Event struct {
	ID        string     `bson:"id"         json:"id"`
	PublicKey string     `bson:"pubkey"     json:"pubkey"`
	CreatedAt int64      `bson:"created_at" json:"created_at"`
	Kind      types.Kind `bson:"kind"       json:"kind"`
	Tags      types.Tags `bson:"tags"       json:"tags"`
	Content   string     `bson:"content"    json:"content"`
	Signature string     `bson:"sig"        json:"sig"`
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

// checkID checks if the user provided id is valid.
func (e *Event) checkID(id [32]byte) bool {
	for i := 0; i < 32; i++ {
		b := hextable[id[i]>>4]
		if b != e.ID[i*2] {
			return false
		}

		b = hextable[id[i]&0x0f]
		if b != e.ID[i*2+1] {
			return false
		}
	}

	return true
}

func (e *Event) checkSig(id [32]byte) bool {
	// turn pubkey hex to byte array.
	pk, err := hex.DecodeString(e.PublicKey)
	if err != nil {
		return false
	}

	// construct the pubkey from byte array.
	pubkey, err := schnorr.ParsePubKey(pk)
	if err != nil {
		return false
	}

	// turn signature hex to byte array.
	s, err := hex.DecodeString(e.Signature)
	if err != nil {
		return false
	}

	// construct signature from byte array.
	sig, err := schnorr.ParseSignature(s)
	if err != nil {
		return false
	}

	// todo::: replace with libsecp256k1 (C++ version).
	return sig.Verify(id[:], pubkey)
}

// IsValid function validats an event Signature and ID.
func (e *Event) IsValid(id [32]byte) bool {
	// make sure the user provided id is valid.
	if !e.checkID(id) {
		return false
	}

	if !e.checkSig(id) {
		return false
	}

	return true
}

// IsProtected checks if ["-"] tag is present, check nip-70 for more.
func (e *Event) IsProtected() bool {
	for _, t := range e.Tags {
		if len(t) != 1 {
			continue
		}

		if t[0] == "-" {
			return true
		}
	}

	return false
}

// String returns and encoded string representation of event e.
func (e *Event) String() string {
	ee, err := e.Encode()
	if err != nil {
		return ""
	}

	return string(ee)
}
