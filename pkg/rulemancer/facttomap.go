package rulemancer

import (
	"errors"
	"fmt"
	"regexp"
)

func genericFactToMap(c *Config, statusItem string, factList string) ([]map[string]string, error) {
	// Pattern to match (statusItem ... anything until matching closing paren
	// We need to manually parse to handle nested parentheses
	var matches [][]string

	escapedItem := regexp.QuoteMeta(statusItem)
	startPattern := regexp.MustCompile(`\(` + escapedItem + `\s+`)

	// Find all start positions
	startMatches := startPattern.FindAllStringIndex(factList, -1)

	if len(startMatches) == 0 {
		return nil, nil
	}

	// For each start position, find the matching closing parenthesis
	for _, startIdx := range startMatches {
		// startIdx[1] is the position right after the pattern match
		pos := startIdx[1]
		depth := 1
		start := pos

		// Find the matching closing paren by counting depth
		for pos < len(factList) && depth > 0 {
			if factList[pos] == '(' {
				depth++
			} else if factList[pos] == ')' {
				depth--
			}
			pos++
		}

		if depth == 0 {
			// Extract the content between parentheses
			content := factList[start : pos-1]
			matches = append(matches, []string{"", content})
		}
	}

	if len(matches) == 0 {
		return nil, errors.New("no matching status items found")
	}

	// Pattern to extract individual key-value pairs: (key value)
	kvPattern := regexp.MustCompile(`\(([^\s]+)\s+([^)]+)\)`)

	var results []map[string]string
	var fieldNames map[string]bool

	// Process each matched statusItem
	for i, match := range matches {
		if len(match) < 2 {
			continue
		}

		// Extract all key-value pairs from this statusItem
		kvPairs := match[1]
		kvMatches := kvPattern.FindAllStringSubmatch(kvPairs, -1)

		itemMap := make(map[string]string)
		for _, kv := range kvMatches {
			if len(kv) >= 3 {
				key := kv[1]
				value := kv[2]
				itemMap[key] = value
			}
		}

		if len(itemMap) == 0 {
			continue
		}

		// Validate that all items have the same fields
		if i == 0 {
			// First item: store field names
			fieldNames = make(map[string]bool)
			for key := range itemMap {
				fieldNames[key] = true
			}
		} else {
			// Subsequent items: check consistency
			if len(itemMap) != len(fieldNames) {
				return nil, fmt.Errorf("inconsistent fields: item %d has %d fields, expected %d", i, len(itemMap), len(fieldNames))
			}
			for key := range itemMap {
				if !fieldNames[key] {
					return nil, fmt.Errorf("inconsistent fields: item %d has unexpected field %q", i, key)
				}
			}
			for expectedKey := range fieldNames {
				if _, exists := itemMap[expectedKey]; !exists {
					return nil, fmt.Errorf("inconsistent fields: item %d missing expected field %q", i, expectedKey)
				}
			}
		}

		results = append(results, itemMap)
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
