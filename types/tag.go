package types

type (
	Tag  []string
	Tags []Tag
)

// ContainsTag checks if a tag with an specific name/value is exist.
func (tags Tags) ContainsTag(name, value string) bool {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] == name && tag[1] == value {
			return true
		}
	}

	return false
}

// ContainsAny checks if event have a tag with given name that its value is equal to one of given values.
func (tags Tags) ContainsAny(name string, values []string) bool {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if "#"+tag[0] != name {
			continue
		}

		if ContainsString(tag[1], values) {
			return true
		}
	}

	return false
}

// GetValue returns the value of first tag with given name.
func (tags Tags) GetValue(name string) string {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] == name {
			return tag[1]
		}
	}

	return ""
}

// GetValues returns all values of all tags with same name in an event.
func (tags Tags) GetValues(name string) []string {
	values := []string{}
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] == name {
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
