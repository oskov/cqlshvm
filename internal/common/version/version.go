package version

import (
	"fmt"
	"regexp"
	"strconv"
)

type Version struct {
	Major int
	Minor int
	Patch int
	RC    int
}

var versionRe = regexp.MustCompile(`^(\d+)\.?(\d+)?\.?(\d+)?(~rc(\d+))?$`)

func Parse(str string) (Version, error) {
	matches := versionRe.FindStringSubmatch(str)
	if len(matches) == 0 {
		return Version{}, fmt.Errorf("invalid version: %s", str)
	}

	parseInt := func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return i
	}

	v := Version{
		RC: -1,
	}
	v.Major = parseInt(matches[1])
	v.Minor = parseInt(matches[2])
	if matches[3] != "" {
		v.Patch = parseInt(matches[3])
	}
	if matches[4] != "" {
		v.RC = parseInt(matches[5])
	}

	return v, nil
}

func (v Version) String() string {
	if v.RC == -1 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d.%d~rc%d", v.Major, v.Minor, v.Patch, v.RC)
}

func (v Version) Eq(other Version) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch && v.RC == other.RC
}

func (v Version) Gt(other Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	if v.Patch != other.Patch {
		return v.Patch > other.Patch
	}
	if v.RC != other.RC {
		if v.RC == -1 || other.RC == -1 {
			return v.RC < other.RC
		}
		return v.RC > other.RC
	}
	return false
}

func (v Version) Gte(other Version) bool {
	return v.Gt(other) || v.Eq(other)
}

func (v Version) Lt(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	if v.Patch != other.Patch {
		return v.Patch < other.Patch
	}
	if v.RC != other.RC {
		if v.RC == -1 || other.RC == -1 {
			return v.RC > other.RC
		}
		return v.RC < other.RC
	}
	return false
}

func (v Version) Lte(other Version) bool {
	return v.Lt(other) || v.Eq(other)
}
