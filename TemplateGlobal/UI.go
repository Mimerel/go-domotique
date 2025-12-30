package TemplateGlobal

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
var hyphenUnderscore = regexp.MustCompile(`[-_]+`)

func GetUIDict() map[string]interface{} {
	return template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("Invalid Dict Call")
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
		},
		"cssClass": func(s string) string {
			// Convert to lowercase and replace non-alphanumeric chars with hyphens
			result := strings.ToLower(s)
			result = nonAlphanumeric.ReplaceAllString(result, "-")
			result = strings.Trim(result, "-")
			return result
		},
		"displayName": func(s string) string {
			// Replace hyphens and underscores with spaces for display
			result := hyphenUnderscore.ReplaceAllString(s, " ")
			return strings.TrimSpace(result)
		},
	}
}
