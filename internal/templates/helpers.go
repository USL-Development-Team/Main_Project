package templates

import (
	"fmt"
	"html/template"
)

// TemplateFunctions returns a map of template helper functions
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"dict":   createDict,
		"slice":  createSlice,
		"add":    addInts,
		"sub":    subtractFloats,
		"mul":    multiplyFloats,
		"printf": sprintf,
		"lt":     lessThan,
		"gt":     greaterThan,
		"eq":     equals,
		"substr": subString,
	}
}

// createDict creates a map from alternating key-value pairs
func createDict(values ...any) map[string]any {
	dict := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		if i+1 < len(values) {
			dict[values[i].(string)] = values[i+1]
		}
	}
	return dict
}

// createSlice creates a slice from provided values
func createSlice(values ...any) []any {
	return values
}

// addInts adds two integers
func addInts(a, b int) int {
	return a + b
}

// subtractFloats subtracts two float64 values
func subtractFloats(a, b float64) float64 {
	return a - b
}

// multiplyFloats multiplies two float64 values
func multiplyFloats(a, b float64) float64 {
	return a * b
}

// sprintf formats a string with provided arguments
func sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// lessThan compares two float64 values
func lessThan(a, b float64) bool {
	return a < b
}

// greaterThan compares two float64 values
func greaterThan(a, b float64) bool {
	return a > b
}

// equals compares two values for equality
func equals(a, b any) bool {
	return a == b
}

// subString extracts a substring from a string
func subString(s string, start, length int) string {
	if start >= len(s) {
		return ""
	}
	end := min(start+length, len(s))
	return s[start:end]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
