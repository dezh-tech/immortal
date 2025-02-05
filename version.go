package immortal

import "fmt"

// These constants follow the semantic versioning 2.0.0 spec.
// see: http://semver.org
var (
	major = 1
	minor = 0
	patch = 0
	meta  = "beta"
)

func StringVersion() string {
	v := fmt.Sprintf("immortal - %d.%d.%d", major, minor, patch)

	if meta != "" {
		v = fmt.Sprintf("%s-%s", v, meta)
	}

	return v
}
