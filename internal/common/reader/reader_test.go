package reader

import (
	"testing"

	"github.com/oskov/cqlshvm/internal/common/version"
)

func TestPrepareRequest(t *testing.T) {
	tests := []struct {
		prefix      string
		expectedURL string
	}{
		{
			prefix:      "test-prefix",
			expectedURL: "https://s3.amazonaws.com/downloads.scylladb.com?delimiter=%2F&prefix=test-prefix",
		},
		{
			prefix:      "",
			expectedURL: "https://s3.amazonaws.com/downloads.scylladb.com?delimiter=%2F&prefix=",
		},
	}

	for _, tt := range tests {
		u, err := prepareURL(s3BucketURL, tt.prefix)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if u != tt.expectedURL {
			t.Errorf("expected URL %v, got %v", tt.expectedURL, u)
		}
	}
}

func TestIsValidPrefix(t *testing.T) {
	tests := []struct {
		prefix string
		params ListParams
		valid  bool
	}{
		{"downloads/scylla-enterprise/relocatable/scylladb-2024.1/", ListParams{}, true},
		{"downloads/scylla-enterprise/relocatable/scylladb-2024.1/",
			ListParams{Gt: &version.Version{
				Major: 2024,
			}},
			true,
		},
		{"downloads/scylla-enterprise/relocatable/scylladb-2024.1/",
			ListParams{Gt: &version.Version{
				Major: 2024,
				Minor: 1,
			}},
			true,
		},
		{"downloads/scylla-enterprise/relocatable/scylladb-2020.1/", ListParams{}, false},
		{"downloads/scylla-enterprise/relocatable/scylladb-branch-2023.1/", ListParams{}, false},
		{"downloads/scylla-enterprise/relocatable/invalid-prefix/", ListParams{}, false},
	}

	for _, test := range tests {
		t.Run(test.prefix, func(t *testing.T) {
			if got := isValidPrefix(test.prefix, test.params); got != test.valid {
				t.Errorf("isValidPrefix(%q) = %v; want %v", test.prefix, got, test.valid)
			}
		})
	}
}
