package main

import (
	"testing"
)

func TestRoundFloat(t *testing.T) {
	tests := []struct {
		value     float64
		precision int
		expected  float64
	}{
		{123.456789, 2, 123.46},
		{123.456789, 3, 123.457},
		{123.456789, 0, 123},
		{123.4, 2, 123.40},
		{123.5, 0, 124},
	}

	for _, test := range tests {
		result := roundFloat(test.value, test.precision)
		if result != test.expected {
			t.Errorf("roundFloat(%f, %d) = %f; expected %f", test.value, test.precision, result, test.expected)
		}
	}
}
