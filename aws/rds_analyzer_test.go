package aws

import (
	"testing"
)

func TestSupportedRDSOfferings(t *testing.T) {
	cases := []struct {
		input    string
		expected ComputeResources
	}{
		{
			input:    "db.t3.micro",
			expected: ComputeResources{CPU: 2, Memory: 1},
		},
		{
			input:    "db.t3.xlarge",
			expected: ComputeResources{CPU: 4, Memory: 16},
		},
		{
			input:    "db.t4g.medium",
			expected: ComputeResources{CPU: 2, Memory: 4},
		},
		{
			input:    "db.m5.large",
			expected: ComputeResources{CPU: 2, Memory: 8},
		},
		{
			input:    "db.r6g.large",
			expected: ComputeResources{CPU: 2, Memory: 8},
		},
	}

	for _, tt := range cases {
		value, err := rdsOfferings(tt.input)
		if err != nil {
			t.Fatalf("input: %s, expect no error, got %v", tt.input, err)
		}
		if *value != tt.expected {
			t.Errorf("Supported Types returned: %d, expected: %d", value, tt.expected)
		}
	}
}

func TestNotSupportedRDSOfferings(t *testing.T) {
	cases := []struct {
		input string
	}{
		{input: "ds23.xlarge"},
		{input: "t3.micro"},
	}

	for _, tt := range cases {
		_, err := rdsOfferings(tt.input)
		if err == nil {
			t.Errorf("Not Supported Types. input: %s, expected: exception", tt.input)
		}
	}
}
