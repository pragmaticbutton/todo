package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceWildcards(t *testing.T) {

	tests := []struct {
		testName string
		name     string
		exp      string
	}{
		{
			testName: "No wildcards",
			name:     "nowildcards",
			exp:      "nowildcards",
		},
		{
			testName: "wildcards 1",
			name:     "word*",
			exp:      "word%",
		},
		{
			testName: "wildcards 2",
			name:     "*word*",
			exp:      "%word%",
		},
		{
			testName: "wildcards 3",
			name:     "*wo*rd*",
			exp:      "%wo%rd%",
		},
		{
			testName: "wildcards 4",
			name:     "*wo**rd*",
			exp:      "%wo%%rd%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			e := replaceWildCards(tt.name)
			assert.Equal(t, tt.exp, e)
		})
	}
}
