package rulemancer

import (
	"strings"
)

// factsSplit splits a string by spaces while preserving quoted strings as single elements.
// Quoted strings are enclosed in double quotes (").
// Multiple spaces are treated as a single separator.
// Example: `prova ciao "uno due tre"` returns ["prova", "ciao", "uno due tre"]
func factsSplit(input string) []string {
	result := []string{}
	var current strings.Builder
	inQuotes := false

	input = strings.TrimSpace(input)

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ', '\t':
			if inQuotes {
				// Inside quotes, keep the space
				current.WriteByte(char)
			} else {
				// Outside quotes, end current token if any
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(char)
		}
	}

	// Add the last token if any
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
