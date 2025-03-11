package reader

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/oskov/cqlshvm/internal/common/version"
)

type ParsedObjectKey struct {
	PrefixVersion version.Version
	ObjectVersion version.Version
}

const cqlshAppName = "scylla-enterprise-cqlsh"

// cqlshRe matches object key
var cqlshRe = regexp.MustCompile(`^.*?scylladb-(\d{4})\.(\d)/?(` + cqlshAppName + `-(\d{4})\.(\d)\.(\d)(~rc(\d))?-).*$`)

// prefixRe matches only prefixes like scylladb-2024.2
var prefixRe = regexp.MustCompile(`^.*?scylladb-(\d{4})\.(\d)/?$`)

func ParseObjectKey(key string) (*ParsedObjectKey, error) {
	matches := prefixRe.FindStringSubmatch(key)
	if len(matches) == 0 || len(matches) < 3 {
		matches = cqlshRe.FindStringSubmatch(key)
		if len(matches) == 0 || len(matches) < 7 {
			return nil, fmt.Errorf("invalid object key: %s", key)
		}
	}

	data := &ParsedObjectKey{
		PrefixVersion: version.Version{
			RC: -1,
		},
		ObjectVersion: version.Version{
			RC: -1,
		},
	}

	var parsingErr error

	parseInt := func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			parsingErr = err
			return 0
		}
		return i
	}

	data.PrefixVersion.Major = parseInt(matches[1])
	data.PrefixVersion.Minor = parseInt(matches[2])
	if len(matches) > 3 && matches[3] != "" {
		data.ObjectVersion.Major = parseInt(matches[4])
		data.ObjectVersion.Minor = parseInt(matches[5])
		data.ObjectVersion.Patch = parseInt(matches[6])
		if matches[7] != "" {
			data.ObjectVersion.RC = parseInt(matches[8])
		}
	}

	return data, parsingErr
}
