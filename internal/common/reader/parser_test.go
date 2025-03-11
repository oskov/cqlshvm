package reader

import (
	"reflect"
	"testing"

	"github.com/oskov/cqlshvm/internal/common/version"
)

func TestParseObjectKey(t *testing.T) {
	tests := []struct {
		input    string
		expected *ParsedObjectKey
		wantErr  bool
	}{
		{
			input: "downloads/scylla-enterprise/relocatable/scylladb-2024.2/scylla-enterprise-cqlsh-2024.2.5-0.20250221.cb9e2a54ae6d.noarch.tar.gz",
			expected: &ParsedObjectKey{
				PrefixVersion: version.Version{Major: 2024, Minor: 2, RC: -1},
				ObjectVersion: version.Version{Major: 2024, Minor: 2, Patch: 5, RC: -1},
			},
			wantErr: false,
		},
		{
			input:    "downloads/scylla-enterprise/relocatable/scylladb-2024.2/",
			expected: &ParsedObjectKey{PrefixVersion: version.Version{Major: 2024, Minor: 2, RC: -1}, ObjectVersion: version.Version{RC: -1}},
			wantErr:  false,
		},
		{
			input: "downloads/scylla-enterprise/relocatable/scylladb-2024.2/scylla-enterprise-cqlsh-2024.2.0~rc2-0.20240904.4c26004e5311.noarch.tar.gz",
			expected: &ParsedObjectKey{
				PrefixVersion: version.Version{Major: 2024, Minor: 2, RC: -1},
				ObjectVersion: version.Version{Major: 2024, Minor: 2, Patch: 0, RC: 2},
			},
			wantErr: false,
		},
		{
			input:   "downloads/scylla-enterprise/relocatable/scylladb-2024.2/scylla-enterprise-jmx-2024.2.0~rc2-0.20240904.4c26004e5311.noarch.tar.gz",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseObjectKey(test.input)
			if !test.wantErr && err != nil {
				t.Fatalf("ParseObjectKey(%q) returned error: %v", test.input, err)
			}
			if test.wantErr && err == nil {
				t.Fatalf("ParseObjectKey(%q) expected error, got nil", test.input)

			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("ParseObjectKey(%q) = %+v; want %+v", test.input, result, test.expected)
			}
		})
	}
}
