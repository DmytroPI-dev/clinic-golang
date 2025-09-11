// /internal/utils/utils.go
package utils

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ToTitle capitalizes the first letter of a string.
func ToTitle(s string) string {
	return cases.Title(language.English).String(s)
}

// Dict creates a map from a list of key-value pairs.
// This is a very useful helper for Go templates.
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call: odd number of arguments")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}