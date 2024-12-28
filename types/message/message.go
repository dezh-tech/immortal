package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/mailru/easyjson"
	jwriter "github.com/mailru/easyjson/jwriter"
	"github.com/tidwall/gjson" // TODO::: remove/replace me!
)

// Message reperesents an NIP-01 message which can be sent to or received by client.
type Message interface {
	Type() string
	DecodeFromJSON([]byte) error
	EncodeToJSON() ([]byte, error)
	String() string
}

// ParseMessage parses the given message from client to a message interface.
func ParseMessage(message []byte) (Message, error) {
	firstComma := bytes.Index(message, []byte{','})
	if firstComma == -1 {
		return nil, errors.New("invalid message: can't find a , in message")
	}
	msgType := message[0:firstComma]

	var e Message
	switch {
	case bytes.Contains(msgType, []byte("EVENT")):
		e = &Event{
			Event: new(event.Event),
		}
	case bytes.Contains(msgType, []byte("REQ")):
		e = &Req{}
	case bytes.Contains(msgType, []byte("AUTH")):
		e = &Auth{}
	case bytes.Contains(msgType, []byte("CLOSE")):
		x := Close("")
		e = &x
	default:
		return nil, errors.New("invalid message type: must be one of REQ, EVENT, AUTH or CLOSE")
	}

	if err := e.DecodeFromJSON(message); err != nil {
		return nil, err
	}

	return e, nil
}

// Event represents a NIP-01 EVENT message.
type Event struct {
	SubscriptionID string
	Event          *event.Event
}

// MakeEvent constructs an EVENT message to be sent to client.
func MakeEvent(id string, e *event.Event) []byte {
	em := Event{
		SubscriptionID: id,
		Event:          e,
	}

	res, err := em.EncodeToJSON()
	if err != nil {
		return []byte{} // TODO::: should we return anything else here?
	}

	return res
}

func (Event) Type() string { return "EVENT" }

func (em Event) String() string {
	v, err := json.Marshal(em)
	if err != nil {
		return ""
	}

	return string(v)
}

func (em *Event) DecodeFromJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	switch {
	case len(arr) >= 2:
		err := easyjson.Unmarshal([]byte(arr[1].Raw), em.Event)
		if err != nil {
			return types.DecodeError{
				Reason: fmt.Sprintf("EVENT message: %s", err.Error()),
			}
		}

		return nil
	default:

		return types.DecodeError{
			Reason: "EVENT messag: no event found.",
		}
	}
}

func (em Event) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["EVENT","` + em.SubscriptionID + `",`)
	em.Event.MarshalEasyJSON(&w)
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("EVENT message: %s", err.Error()),
		}
	}

	return res, nil
}

// Req represents a NIP-01 REQ message.
type Req struct {
	SubscriptionID string
	filter.Filters
}

func (Req) Type() string { return "REQ" }

func (rm *Req) DecodeFromJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 3 {
		return types.DecodeError{
			Reason: "REQ message: missing filters.",
		}
	}
	rm.SubscriptionID = arr[1].Str
	rm.Filters = make(filter.Filters, len(arr)-2)
	f := 0
	for i := 2; i < len(arr); i++ {
		if err := easyjson.Unmarshal([]byte(arr[i].Raw), &rm.Filters[f]); err != nil {
			return types.DecodeError{
				Reason: fmt.Sprintf("REQ message: %s", err.Error()),
			}
		}
		f++
	}

	return nil
}

func (rm Req) EncodeToJSON() ([]byte, error) {
	return nil, nil
}

// Notice reperesents a NIP-01 NOTICE message.
type Notice string

func MakeNotice(msg string) []byte {
	res, err := Notice(msg).EncodeToJSON()
	if err != nil {
		return []byte{} // TODO::: should we return anything else here?
	}

	return res
}

func (Notice) Type() string { return "NOTICE" }
func (nm Notice) String() string {
	v, err := json.Marshal(nm)
	if err != nil {
		return ""
	}

	return string(v)
}

func (nm *Notice) DecodeFromJSON(_ []byte) error {
	return nil
}

func (nm Notice) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["NOTICE",`)
	w.Raw(json.Marshal(string(nm)))
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("NOTICE message: %s", err.Error()),
		}
	}

	return res, nil
}

// EOSE reperesents a NIP-01 EOSE message.
type EOSE string

func MakeEOSE(sID string) []byte {
	res, err := EOSE(sID).EncodeToJSON()
	if err != nil {
		return []byte{} // TODO::: should we return anything else here?
	}

	return res
}

func (EOSE) Type() string { return "EOSE" }
func (em EOSE) String() string {
	v, err := json.Marshal(em)
	if err != nil {
		return ""
	}

	return string(v)
}

func (em *EOSE) DecodeFromJSON(_ []byte) error {
	return nil
}

func (em EOSE) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["EOSE",`)
	w.Raw(json.Marshal(string(em)))
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("EOSE message: %s", err.Error()),
		}
	}

	return res, nil
}

// Close reperesents a NIP-01 CLOSE message.
type Close string

func (Close) Type() string { return "CLOSE" }
func (cm Close) String() string {
	return string(cm)
}

func (cm *Close) DecodeFromJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	switch len(arr) {
	case 2:
		*cm = Close(arr[1].Str)

		return nil
	default:

		return types.DecodeError{
			Reason: "CLOSE message: subscription ID missed.",
		}
	}
}

func (cm Close) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["CLOSE",`)
	w.Raw(json.Marshal(string(cm)))
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("CLOSE message: %s", err.Error()),
		}
	}

	return res, nil
}

// Closed reperesents a NIP-01 CLOSED message.
type Closed struct {
	SubscriptionID string
	Reason         string
}

// MakeClosed constructs a CLOSED message to be sent to client.
func MakeClosed(id, reason string) []byte {
	cm := Closed{
		SubscriptionID: id,
		Reason:         reason,
	}

	res, err := cm.EncodeToJSON()
	if err != nil {
		return []byte{}
	}

	return res
}

func (Closed) Label() string { return "CLOSED" }
func (cm Closed) String() string {
	v, err := json.Marshal(cm)
	if err != nil {
		return ""
	}

	return string(v)
}

func (cm *Closed) DecodeFromJSON(_ []byte) error {
	return nil
}

func (cm Closed) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["CLOSED",`)
	w.Raw(json.Marshal(cm.SubscriptionID))
	w.RawString(`,`)
	w.Raw(json.Marshal(cm.Reason))
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("CLOSED message: %s", err.Error()),
		}
	}

	return res, nil
}

// OK reperesents a NIP-01 OK message.
type OK struct {
	OK      bool
	EventID string
	Reason  string
}

// MakeOK constructs a NIP-01 OK message to be sent to the client.
func MakeOK(ok bool, eid, reason string) []byte {
	om := OK{
		OK:      ok,
		EventID: eid,
		Reason:  reason,
	}

	res, err := om.EncodeToJSON()
	if err != nil {
		return []byte{} // TODO::: should we return anything else here?
	}

	return res
}

func (OK) Type() string { return "OK" }
func (om OK) String() string {
	v, err := json.Marshal(om)
	if err != nil {
		return ""
	}

	return string(v)
}

func (om *OK) DecodeFromJSON(_ []byte) error {
	return nil
}

func (om OK) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["OK","` + om.EventID + `",`)
	ok := "false"
	if om.OK {
		ok = "true"
	}
	w.RawString(ok)
	w.RawString(`,`)
	w.Raw(json.Marshal(om.Reason))
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("OK message: %s", err.Error()),
		}
	}

	return res, nil
}

// Auth reperesents a NIP-01 AUTH message.
type Auth struct {
	Challenge string
	Event     event.Event
}

// MakeAuth constructs a NIP-01 OK message to be sent to the client.
func MakeAuth(challenge string) []byte {
	om := Auth{
		Challenge: challenge,
	}

	res, err := om.EncodeToJSON()
	if err != nil {
		return []byte{} // TODO::: should we return anything else here?
	}

	return res
}

func (Auth) Type() string { return "AUTH" }

func (am Auth) String() string {
	v, err := json.Marshal(am)
	if err != nil {
		return ""
	}

	return string(v)
}

func (am *Auth) DecodeFromJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 2 {
		return types.DecodeError{
			Reason: "AUTH message: missing fields.",
		}
	}

	if arr[1].IsObject() {
		err := easyjson.Unmarshal([]byte(arr[1].Raw), &am.Event)
		if err != nil {
			return types.DecodeError{
				Reason: fmt.Sprintf("AUTH message: %s", err.Error()),
			}
		}

		return nil
	}

	return types.DecodeError{
		Reason: "AUTH is not valid.",
	}
}

func (am *Auth) EncodeToJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["AUTH",`)
	w.Raw(json.Marshal(am.Challenge))
	w.RawString(`]`)

	res, err := w.BuildBytes()
	if err != nil {
		return nil, types.EncodeError{
			Reason: fmt.Sprintf("AUTH message: %s", err.Error()),
		}
	}

	return res, nil
}
