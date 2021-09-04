package parser

import (
	"fmt"
	"testing"
)

func TestDigitalFileSizeString(t *testing.T) {
	testData := []struct {
		input    DigitalFileSize
		expected string
	}{
		{
			DigitalFileSize(11620000000),
			"10.82 GiB",
		},
		{
			DigitalFileSize(4150000000),
			"3.86 GiB",
		},
		{
			DigitalFileSize(774070000),
			"738.21 MiB",
		},
		{
			DigitalFileSize(700),
			"700 Bytes",
		},
	}

	for _, d := range testData {
		result := fmt.Sprintf("%s", d.input)
		if result != d.expected {
			t.Errorf("input: %d\nexpected: %s\ngot: %s", d.input, d.expected, result)
		}
	}
}
