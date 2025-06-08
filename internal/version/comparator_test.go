package version

import "testing"

func TestGoVersionComparator_Compare(t *testing.T) {
	comparator := NewSemanticVersionComparator()

	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "v1 greater than v2 major version",
			v1:       "2.0",
			v2:       "1.23",
			expected: 1,
		},
		{
			name:     "v1 greater than v2 minor version",
			v1:       "1.23",
			v2:       "1.22",
			expected: 1,
		},
		{
			name:     "v1 less than v2 major version",
			v1:       "1.23",
			v2:       "2.0",
			expected: -1,
		},
		{
			name:     "v1 less than v2 minor version",
			v1:       "1.21",
			v2:       "1.22",
			expected: -1,
		},
		{
			name:     "v1 equal to v2",
			v1:       "1.22",
			v2:       "1.22",
			expected: 0,
		},
		{
			name:     "v1 has more parts",
			v1:       "1.22.1",
			v2:       "1.22",
			expected: 1,
		},
		{
			name:     "v2 has more parts",
			v1:       "1.22",
			v2:       "1.22.1",
			expected: -1,
		},
		{
			name:     "already prefixed with go",
			v1:       "go1.23",
			v2:       "go1.22",
			expected: 1,
		},
		{
			name:     "mixed prefixed and non-prefixed",
			v1:       "go1.23",
			v2:       "1.22",
			expected: 1,
		},
		{
			name:     "real Go versions chronological order",
			v1:       "1.21",
			v2:       "1.20",
			expected: 1,
		},
		{
			name:     "patch versions",
			v1:       "1.22.3",
			v2:       "1.22.1",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := comparator.Compare(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("Compare(%s, %s) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestNormalizeGoVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple version",
			input:    "1.22",
			expected: "go1.22",
		},
		{
			name:     "patch version",
			input:    "1.22.1",
			expected: "go1.22.1",
		},
		{
			name:     "already prefixed",
			input:    "go1.22",
			expected: "go1.22",
		},
		{
			name:     "empty version",
			input:    "",
			expected: "go1.0",
		},
		{
			name:     "version with rc",
			input:    "1.23rc1",
			expected: "go1.23rc1",
		},
		{
			name:     "single digit",
			input:    "2",
			expected: "go2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeGoVersion(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeGoVersion(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
