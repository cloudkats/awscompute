package aws

import (
	"testing"
)

func TestSupportedRedshiftType(t *testing.T) {
	cases := []struct {
		input    string
		expected ComputeResources
	}{
		{
			input:    "ds2.xlarge",
			expected: ComputeResources{CPU: 4, Memory: 31},
		},
		{
			input:    "dc1.8xlarge",
			expected: ComputeResources{CPU: 32, Memory: 244},
		},
		{
			input:    "ra3.4xlarge",
			expected: ComputeResources{CPU: 12, Memory: 96},
		},
	}

	for _, tt := range cases {
		value, err := redshiftOfferings(tt.input)
		if err != nil {
			t.Fatalf("input: %s, expect no error, got %v", tt.input, err)
		}
		if *value != tt.expected {
			t.Errorf("Supported Types returned: %d, expected: %d", value, tt.expected)
		}
	}
}

func TestNotSupportedRedshiftType(t *testing.T) {
	cases := []struct {
		input string
	}{
		{input: "ds23.xlarge"},
		{input: "t3.micro"},
		{input: "t3.micro"},
	}

	for _, tt := range cases {
		_, err := redshiftOfferings(tt.input)
		if err == nil {
			t.Errorf("Not Supported Types. input: %s, expected: exception", tt.input)
		}
	}
}
