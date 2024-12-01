package data

import (
	"testing"

	"greenlight.honganhpham.net/internal/assert"
)

func TestRuntime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		runtime  Runtime
		expected string
	}{
		{
			name:     "Valid runtime marshaling",
			runtime:  120,
			expected: `"120 mins"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonValue, err := test.runtime.MarshalJSON()
			if err != nil {
				assert.NilError(t, err)
			}
			if string(jsonValue) != test.expected {
				assert.Equal(t, string(jsonValue), test.expected)

			}
		})
	}
}

func TestRuntime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Runtime
		expectError bool
	}{
		{
			name:        "Valid runtime unmarshaling",
			input:       `"120 mins"`,
			expected:    120,
			expectError: false,
		},
		{
			name:        "Invalid format - no quotes",
			input:       `120 mins`,
			expectError: true,
		},
		{
			name:        "Invalid format - wrong suffix",
			input:       `"120 minutes"`,
			expectError: true,
		},
		{
			name:        "Invalid format - no space",
			input:       `"120mins"`,
			expectError: true,
		},
		{
			name:        "Invalid number",
			input:       `"abc mins"`,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var runtime Runtime
			err := runtime.UnmarshalJSON([]byte(test.input))

			if test.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !test.expectError {
				if err != nil {
					assert.NilError(t, err)
				}
				if runtime != test.expected {
					assert.Equal(t, runtime, test.expected)
				}
			}
		})
	}
}
