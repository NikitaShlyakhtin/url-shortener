package data

import (
	"testing"
)

func TestGenerateUrlSuffix(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "Test with length 5",
			length: 5,
		},
		{
			name:   "Test with length 10",
			length: 10,
		},
		{
			name:   "Test with length 15",
			length: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix, err := generateUrlSuffix(tt.length)
			if err != nil {
				t.Errorf("generateUrlSuffix() error = %v", err)
				return
			}

			if len(suffix) != tt.length {
				t.Errorf("generateUrlSuffix() = %v, want %v", len(suffix), tt.length)
			}
		})
	}
}
