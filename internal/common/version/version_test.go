package version

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
		hasError bool
	}{
		{"1.0.0", Version{1, 0, 0, -1}, false},
		{"2.1.3", Version{2, 1, 3, -1}, false},
		{"3.2.1~rc1", Version{3, 2, 1, 1}, false},
		{"4.0.0~rc2", Version{4, 0, 0, 2}, false},
		{"5.0", Version{5, 0, 0, -1}, false},
		{"5", Version{5, 0, 0, -1}, false},
		{"invalid", Version{}, true},
		{"6.0.0~rc", Version{}, true}, // no rc number
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Parse(test.input)
			if (err != nil) != test.hasError {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			}
			if !test.hasError && result != test.expected {
				t.Errorf("Parse(%q) = %v, want %v", test.input, result, test.expected)
			}
		})
	}
}

func TestVersionGt(t *testing.T) {
	tests := []struct {
		v1, v2   string
		expected bool
	}{
		{"1.0.0", "1.0.0~rc1", true},
		{"1.0.0~rc2", "1.0.0~rc1", true},
		{"1.0.1", "1.0.0~rc1", true},
		{"2.0.0", "1.0.0", true},
		{"1.1.0", "1.0.0", true},
		{"1.0.0~rc1", "1.0.0", false},
		{"1.0.0", "1.0.0", false},
	}

	for _, test := range tests {
		t.Run(test.v1+" > "+test.v2, func(t *testing.T) {
			v1, err := Parse(test.v1)
			if err != nil {
				t.Fatalf("Gt(%q, %q) error = %v", test.v1, test.v2, err)
			}
			v2, err := Parse(test.v2)
			if err != nil {
				t.Fatalf("Gt(%q, %q) error = %v", test.v1, test.v2, err)
			}
			got := v1.Gt(v2)
			if got != test.expected {
				t.Errorf("Gt(%q, %q) = %v; want %v", test.v1, test.v2, got, test.expected)
			}
		})
	}
}

func TestVersionLt(t *testing.T) {
	tests := []struct {
		v1, v2   string
		expected bool
	}{
		{"1.0.0~rc1", "1.0.0", true},
		{"1.0.0~rc1", "1.0.1", true},
		{"1.0.0~rc1", "1.0.0~rc2", true},
		{"1.0.0", "2.0.0", true},
		{"1.0.0", "1.1.0", true},
		{"1.0.0", "1.0.0~rc1", false},
		{"1.0.0", "1.0.0", false},
	}

	for _, test := range tests {
		t.Run(test.v1+" < "+test.v2, func(t *testing.T) {
			v1, err := Parse(test.v1)
			if err != nil {
				t.Fatalf("Lt(%q, %q) error = %v", test.v1, test.v2, err)
			}
			v2, err := Parse(test.v2)
			if err != nil {
				t.Fatalf("Lt(%q, %q) error = %v", test.v1, test.v2, err)
			}
			got := v1.Lt(v2)
			if err != nil {
				t.Fatalf("Lt(%q, %q) error = %v", test.v1, test.v2, err)
			}
			if got != test.expected {
				t.Errorf("Lt(%q, %q) = %v; want %v", test.v1, test.v2, got, test.expected)
			}
		})
	}
}
