package envelope

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/mailru/easyjson"
	jwriter "github.com/mailru/easyjson/jwriter"
	"github.com/tidwall/gjson" // TODO::: remove/replace me!
)

func ParseMessage(message []byte) Envelope {
	firstComma := bytes.Index(message, []byte{','})
	if firstComma == -1 {
		return nil
	}
	label := message[0:firstComma]

	var v Envelope
	switch {
	case bytes.Contains(label, []byte("EVENT")):
		v = &EventEnvelope{}
	case bytes.Contains(label, []byte("REQ")):
		v = &ReqEnvelope{}
	case bytes.Contains(label, []byte("COUNT")):
		v = &CountEnvelope{}
	case bytes.Contains(label, []byte("NOTICE")):
		x := NoticeEnvelope("")
		v = &x
	case bytes.Contains(label, []byte("EOSE")):
		x := EOSEEnvelope("")
		v = &x
	case bytes.Contains(label, []byte("OK")):
		v = &OKEnvelope{}
	case bytes.Contains(label, []byte("AUTH")):
		v = &AuthEnvelope{}
	case bytes.Contains(label, []byte("CLOSED")):
		v = &ClosedEnvelope{}
	case bytes.Contains(label, []byte("CLOSE")):
		x := CloseEnvelope("")
		v = &x
	default:
		return nil
	}

	if err := v.UnmarshalJSON(message); err != nil {
		return nil
	}

	return v
}

type Envelope interface {
	Label() string
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)
	String() string
}

var (
	_ Envelope = (*EventEnvelope)(nil)
	_ Envelope = (*ReqEnvelope)(nil)
	_ Envelope = (*CountEnvelope)(nil)
	_ Envelope = (*NoticeEnvelope)(nil)
	_ Envelope = (*EOSEEnvelope)(nil)
	_ Envelope = (*CloseEnvelope)(nil)
	_ Envelope = (*OKEnvelope)(nil)
	_ Envelope = (*AuthEnvelope)(nil)
)

type EventEnvelope struct {
	SubscriptionID string
	Event          *event.Event
}

func (EventEnvelope) Label() string { return "EVENT" }

func (ee EventEnvelope) String() string {
	v, err := json.Marshal(ee)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ee *EventEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	switch len(arr) {
	case 2:

		return easyjson.Unmarshal([]byte(arr[1].Raw), ee.Event)
	case 3:
		ee.SubscriptionID = arr[1].Str

		return easyjson.Unmarshal([]byte(arr[2].Raw), ee.Event)
	default:

		return fmt.Errorf("failed to decode EVENT envelope")
	}
}

func (ee EventEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["EVENT",`)

	if ee.SubscriptionID != "" {
		w.RawString(`"` + ee.SubscriptionID + `",`)
	}

	ee.Event.MarshalEasyJSON(&w)
	w.RawString(`]`)

	return w.BuildBytes()
}

type ReqEnvelope struct {
	SubscriptionID string
	filter.Filters
}

func (ReqEnvelope) Label() string { return "REQ" }

func (re *ReqEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 3 {
		return fmt.Errorf("failed to decode REQ envelope: missing filters")
	}
	re.SubscriptionID = arr[1].Str
	re.Filters = make(filter.Filters, len(arr)-2)
	f := 0
	for i := 2; i < len(arr); i++ {
		if err := easyjson.Unmarshal([]byte(arr[i].Raw), &re.Filters[f]); err != nil {
			return fmt.Errorf("%w -- on filter %d", err, f)
		}
		f++
	}

	return nil
}

func (re ReqEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["REQ",`)
	w.RawString(`"` + re.SubscriptionID + `"`)
	for _, filter := range re.Filters {
		w.RawString(`,`)
		filter.MarshalEasyJSON(&w)
	}
	w.RawString(`]`)

	return w.BuildBytes()
}

type CountEnvelope struct {
	SubscriptionID string
	Filters        []*filter.Filter
	Count          *int64
}

func (CountEnvelope) Label() string { return "COUNT" }
func (ce CountEnvelope) String() string {
	v, err := json.Marshal(ce)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ce *CountEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 3 {
		return fmt.Errorf("failed to decode COUNT envelope: missing filters")
	}
	ce.SubscriptionID = arr[1].Str

	if len(arr) < 3 {
		return fmt.Errorf("COUNT array must have at least 3 items")
	}

	var countResult struct {
		Count *int64 `json:"count"`
	}
	if err := json.Unmarshal([]byte(arr[2].Raw), &countResult); err == nil && countResult.Count != nil {
		ce.Count = countResult.Count

		return nil
	}

	ce.Filters = make([]*filter.Filter, len(arr)-2)
	f := 0
	for i := 2; i < len(arr); i++ {
		item := []byte(arr[i].Raw)

		if err := easyjson.Unmarshal(item, ce.Filters[f]); err != nil {
			return fmt.Errorf("%w -- on filter %d", err, f)
		}

		f++
	}

	return nil
}

func (ce CountEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["COUNT",`)
	w.RawString(`"` + ce.SubscriptionID + `"`)
	if ce.Count != nil {
		w.RawString(`,{"count":`)
		w.RawString(strconv.FormatInt(*ce.Count, 10))
		w.RawString(`}`)
	} else {
		for _, filter := range ce.Filters {
			w.RawString(`,`)
			filter.MarshalEasyJSON(&w)
		}
	}
	w.RawString(`]`)

	return w.BuildBytes()
}

type NoticeEnvelope string

func (NoticeEnvelope) Label() string { return "NOTICE" }
func (ne NoticeEnvelope) String() string {
	v, err := json.Marshal(ne)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ne *NoticeEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 2 {
		return fmt.Errorf("failed to decode NOTICE envelope")
	}
	*ne = NoticeEnvelope(arr[1].Str)

	return nil
}

func (ne NoticeEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["NOTICE",`)
	w.Raw(json.Marshal(string(ne)))
	w.RawString(`]`)

	return w.BuildBytes()
}

type EOSEEnvelope string

func (EOSEEnvelope) Label() string { return "EOSE" }
func (ee EOSEEnvelope) String() string {
	v, err := json.Marshal(ee)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ee *EOSEEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 2 {
		return fmt.Errorf("failed to decode EOSE envelope")
	}
	*ee = EOSEEnvelope(arr[1].Str)

	return nil
}

func (ee EOSEEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["EOSE",`)
	w.Raw(json.Marshal(string(ee)))
	w.RawString(`]`)

	return w.BuildBytes()
}

type CloseEnvelope string

func (CloseEnvelope) Label() string { return "CLOSE" }
func (ce CloseEnvelope) String() string {
	v, err := json.Marshal(ce)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ce *CloseEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	switch len(arr) {
	case 2:
		*ce = CloseEnvelope(arr[1].Str)

		return nil
	default:

		return fmt.Errorf("failed to decode CLOSE envelope")
	}
}

func (ce CloseEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["CLOSE",`)
	w.Raw(json.Marshal(string(ce)))
	w.RawString(`]`)

	return w.BuildBytes()
}

type ClosedEnvelope struct {
	SubscriptionID string
	Reason         string
}

func (ClosedEnvelope) Label() string { return "CLOSED" }
func (ce ClosedEnvelope) String() string {
	v, err := json.Marshal(ce)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ce *ClosedEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	switch len(arr) {
	case 3:
		*ce = ClosedEnvelope{arr[1].Str, arr[2].Str}

		return nil
	default:

		return fmt.Errorf("failed to decode CLOSED envelope")
	}
}

func (ce ClosedEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["CLOSED",`)
	w.Raw(json.Marshal(ce.SubscriptionID))
	w.RawString(`,`)
	w.Raw(json.Marshal(ce.Reason))
	w.RawString(`]`)

	return w.BuildBytes()
}

type OKEnvelope struct {
	EventID string
	OK      bool
	Reason  string
}

func (OKEnvelope) Label() string { return "OK" }
func (oe OKEnvelope) String() string {
	v, err := json.Marshal(oe)
	if err != nil {
		return ""
	}

	return string(v)
}

func (oe *OKEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 4 {
		return fmt.Errorf("failed to decode OK envelope: missing fields")
	}
	oe.EventID = arr[1].Str
	oe.OK = arr[2].Raw == "true"
	oe.Reason = arr[3].Str

	return nil
}

func (oe OKEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["OK",`)
	w.RawString(`"` + oe.EventID + `",`)
	ok := "false"
	if oe.OK {
		ok = "true"
	}
	w.RawString(ok)
	w.RawString(`,`)
	w.Raw(json.Marshal(oe.Reason))
	w.RawString(`]`)

	return w.BuildBytes()
}

type AuthEnvelope struct {
	Challenge *string
	Event     *event.Event
}

func (AuthEnvelope) Label() string { return "AUTH" }

func (ae AuthEnvelope) String() string {
	v, err := json.Marshal(ae)
	if err != nil {
		return ""
	}

	return string(v)
}

func (ae *AuthEnvelope) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	arr := r.Array()
	if len(arr) < 2 {
		return fmt.Errorf("failed to decode Auth envelope: missing fields")
	}

	if arr[1].IsObject() {
		return easyjson.Unmarshal([]byte(arr[1].Raw), ae.Event)
	}

	ae.Challenge = &arr[1].Str

	return nil
}

func (ae AuthEnvelope) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	w.RawString(`["AUTH",`)
	if ae.Challenge != nil {
		w.Raw(json.Marshal(*ae.Challenge))
	} else {
		ae.Event.MarshalEasyJSON(&w)
	}

	w.RawString(`]`)

	return w.BuildBytes()
}
