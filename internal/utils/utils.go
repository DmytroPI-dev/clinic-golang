// /internal/utils/utils.go
package utils

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Title capitalizes the first letter of a string.
func Title(s string) string {
	return cases.Title(language.English).String(s)
}

// Dict creates a map from a list of key-value pairs.
// This is a very useful helper for Go templates.
func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call: odd number of arguments")
	}
	Dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		Dict[key] = values[i+1]
	}
	return Dict, nil
}