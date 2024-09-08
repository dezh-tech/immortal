package types

type Tag []string

// Marshal Tag. Used for Serialization so string escaping should be as in RFC8259.
func (tag Tag) MarshalTo(dst []byte) []byte {
	dst = append(dst, '[')
	for i, s := range tag {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = EscapeString(dst, s)
	}
	dst = append(dst, ']')

	return dst
}

// MarshalTo appends the JSON encoded byte of Tags as [][]string to dst.
// String escaping is as described in RFC8259.
func MarshalTo(tags []Tag, dst []byte) []byte {
	dst = append(dst, '[')
	for i, tag := range tags {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = tag.MarshalTo(dst)
	}
	dst = append(dst, ']')

	return dst
}
