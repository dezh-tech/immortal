package filter

import (
	json "encoding/json"

	"github.com/dezh-tech/immortal/types"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning.
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson4d398eaaDecodeGithubComDezhTechImmortalTypesFilter(in *jlexer.Lexer, out *Filter) { //nolint
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()

		return
	}
	out.Tags = make(map[string]types.Tag)
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
		case "ids":
			if in.IsNull() {
				in.Skip()
				out.IDs = nil
			} else {
				in.Delim('[')
				if out.IDs == nil {
					if !in.IsDelim(']') {
						out.IDs = make([]string, 0, 20)
					} else {
						out.IDs = []string{}
					}
				} else {
					out.IDs = (out.IDs)[:0]
				}
				for !in.IsDelim(']') {
					v1 := in.String()
					out.IDs = append(out.IDs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "kinds":
			if in.IsNull() {
				in.Skip()
				out.Kinds = nil
			} else {
				in.Delim('[')
				if out.Kinds == nil {
					if !in.IsDelim(']') {
						out.Kinds = make([]types.Kind, 0, 8)
					} else {
						out.Kinds = []types.Kind{}
					}
				} else {
					out.Kinds = (out.Kinds)[:0]
				}
				for !in.IsDelim(']') {
					v2 := types.Kind(in.Int16())
					out.Kinds = append(out.Kinds, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "authors":
			if in.IsNull() {
				in.Skip()
				out.Authors = nil
			} else {
				in.Delim('[')
				if out.Authors == nil {
					if !in.IsDelim(']') {
						out.Authors = make([]string, 0, 40)
					} else {
						out.Authors = []string{}
					}
				} else {
					out.Authors = (out.Authors)[:0]
				}
				for !in.IsDelim(']') {
					v3 := in.String()
					out.Authors = append(out.Authors, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "since":
			if in.IsNull() {
				in.Skip()
				out.Since = 0
			} else {
				if out.Since == 0 {
					out.Since = *new(int64) //nolint
				}
				out.Since = in.Int64()
			}
		case "until":
			if in.IsNull() {
				in.Skip()
				out.Until = 0
			} else {
				if out.Until == 0 {
					out.Until = *new(int64) //nolint
				}
				out.Until = in.Int64()
			}
		case "limit":
			out.Limit = in.Uint16()
		case "search":
			out.Search = in.String()
		default:
			if len(key) > 1 && key[0] == '#' {
				tagValues := make([]string, 0, 40)
				if !in.IsNull() {
					in.Delim('[')
					if out.Authors == nil {
						if !in.IsDelim(']') {
							tagValues = make([]string, 0, 4)
						} else {
							tagValues = []string{}
						}
					} else {
						tagValues = (tagValues)[:0]
					}
					for !in.IsDelim(']') {
						v3 := in.String()
						tagValues = append(tagValues, v3)
						in.WantComma()
					}
					in.Delim(']')
				}
				out.Tags[key[1:]] = tagValues
			} else {
				in.SkipRecursive()
			}
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjson4d398eaaEncodeGithubComDezhTechImmortalTypesFilter(out *jwriter.Writer, in Filter) { //nolint
	out.RawByte('{')
	first := true
	_ = first
	if len(in.IDs) != 0 {
		const prefix string = ",\"ids\":"
		first = false
		out.RawString(prefix[1:])

		out.RawByte('[')
		for v4, v5 := range in.IDs {
			if v4 > 0 {
				out.RawByte(',')
			}
			out.String(v5)
		}
		out.RawByte(']')
	}
	if len(in.Kinds) != 0 {
		const prefix string = ",\"kinds\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}

		out.RawByte('[')
		for v6, v7 := range in.Kinds {
			if v6 > 0 {
				out.RawByte(',')
			}
			out.Int(int(v7))
		}
		out.RawByte(']')
	}
	if len(in.Authors) != 0 {
		const prefix string = ",\"authors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawByte('[')
		for v8, v9 := range in.Authors {
			if v8 > 0 {
				out.RawByte(',')
			}
			out.String(v9)
		}
		out.RawByte(']')
	}
	if in.Since != 0 {
		const prefix string = ",\"since\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(in.Since)
	}
	if in.Until != 0 {
		const prefix string = ",\"until\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(in.Until)
	}
	if in.Limit != 0 || in.Limit == 0 {
		const prefix string = ",\"limit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Limit))
	}
	if in.Search != "" {
		const prefix string = ",\"search\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(in.Search)
	}
	for tag, values := range in.Tags {
		const prefix string = ",\"authors\":" //nolint
		if first {
			first = false
			out.RawString("\"#" + tag + "\":")
		} else {
			out.RawString(",\"#" + tag + "\":")
		}
		out.RawByte('[')
		for i, v := range values {
			if i > 0 {
				out.RawByte(',')
			}
			out.String(v)
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface.
func (f Filter) MarshalJSON() ([]byte, error) { //nolint
	w := jwriter.Writer{}
	easyjson4d398eaaEncodeGithubComDezhTechImmortalTypesFilter(&w, f)

	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface.
func (f Filter) MarshalEasyJSON(w *jwriter.Writer) { //nolint
	easyjson4d398eaaEncodeGithubComDezhTechImmortalTypesFilter(w, f)
}

// UnmarshalJSON supports json.Unmarshaler interface.
func (f *Filter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4d398eaaDecodeGithubComDezhTechImmortalTypesFilter(&r, f)

	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface.
func (f *Filter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4d398eaaDecodeGithubComDezhTechImmortalTypesFilter(l, f)
}
