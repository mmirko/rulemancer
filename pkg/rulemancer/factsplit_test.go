package rulemancer

import (
	"reflect"
	"testing"
)

func TestFactsSplit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple words",
			input:    "prova ciao",
			expected: []string{"prova", "ciao"},
		},
		{
			name:     "quoted string",
			input:    `prova ciao "uno due tre"`,
			expected: []string{"prova", "ciao", "uno due tre"},
		},
		{
			name:     "multiple spaces",
			input:    "prova   ciao    test",
			expected: []string{"prova", "ciao", "test"},
		},
		{
			name:     "quoted string with multiple spaces inside",
			input:    `word "multiple   spaces   here" end`,
			expected: []string{"word", "multiple   spaces   here", "end"},
		},
		{
			name:     "multiple quoted strings",
			input:    `"first quote" middle "second quote"`,
			expected: []string{"first quote", "middle", "second quote"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "single word",
			input:    "word",
			expected: []string{"word"},
		},
		{
			name:     "single quoted string",
			input:    `"quoted string"`,
			expected: []string{"quoted string"},
		},
		{
			name:     "quoted at start",
			input:    `"start quote" middle end`,
			expected: []string{"start quote", "middle", "end"},
		},
		{
			name:     "quoted at end",
			input:    `start middle "end quote"`,
			expected: []string{"start", "middle", "end quote"},
		},
		{
			name:     "tabs as separators",
			input:    "word1	word2	word3",
			expected: []string{"word1", "word2", "word3"},
		},
		{
			name:     "mixed spaces and tabs",
			input:    "word1  	word2	  word3",
			expected: []string{"word1", "word2", "word3"},
		},
		{
			name:     "consecutive quoted strings",
			input:    `"first" "second" "third"`,
			expected: []string{"first", "second", "third"},
		},
		{
			name:     "leading and trailing spaces",
			input:    "  word1 word2  ",
			expected: []string{"word1", "word2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := factsSplit(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("factsSplit(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
