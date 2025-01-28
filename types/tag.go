package types

type (
	Tag  []string
	Tags []Tag
)

func (tags Tags) ContainsAny(tagName string, values []string) bool {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if "#"+tag[0] != tagName {
			continue
		}

		if ContainsString(tag[1], values) {
			return true
		}
	}

	return false
}

func (tags Tags) GetValue(tagName string) string {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] == tagName {
			return tag[1]
		}
	}

	return ""
}

func (tags Tags) GetValues(tagName string) []string {
	values := []string{}
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] == tagName {
			values = append(values, tag[1])
		}
	}

	return values
}

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
