// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package event

import (
	json "encoding/json"
	types "github.com/dezh-tech/immortal/types"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF642ad3eDecodeGithubComDezhTechImmortalTypesEvent(in *jlexer.Lexer, out *Event) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "pubkey":
			out.PublicKey = string(in.String())
		case "created_at":
			out.CreatedAt = int64(in.Int64())
		case "kind":
			out.Kind = types.Kind(in.Uint16())
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]types.Tag, 0, 2)
					} else {
						out.Tags = []types.Tag{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v1 types.Tag
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						in.Delim('[')
						if v1 == nil {
							if !in.IsDelim(']') {
								v1 = make(types.Tag, 0, 4)
							} else {
								v1 = types.Tag{}
							}
						} else {
							v1 = (v1)[:0]
						}
						for !in.IsDelim(']') {
							var v2 string
							v2 = string(in.String())
							v1 = append(v1, v2)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.Tags = append(out.Tags, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "content":
			out.Content = string(in.String())
		case "sig":
			out.Signature = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeGithubComDezhTechImmortalTypesEvent(out *jwriter.Writer, in Event) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"pubkey\":"
		out.RawString(prefix)
		out.String(string(in.PublicKey))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Int64(int64(in.CreatedAt))
	}
	{
		const prefix string = ",\"kind\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Kind))
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Tags {
				if v3 > 0 {
					out.RawByte(',')
				}
				if v4 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v5, v6 := range v4 {
						if v5 > 0 {
							out.RawByte(',')
						}
						out.String(string(v6))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		out.String(string(in.Content))
	}
	{
		const prefix string = ",\"sig\":"
		out.RawString(prefix)
		out.String(string(in.Signature))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Event) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeGithubComDezhTechImmortalTypesEvent(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeGithubComDezhTechImmortalTypesEvent(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeGithubComDezhTechImmortalTypesEvent(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeGithubComDezhTechImmortalTypesEvent(l, v)
}
