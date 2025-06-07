package version

import "testing"

func TestSemanticVersionComparator_Compare(t *testing.T) {
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