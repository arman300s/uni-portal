package contracts

import "strings"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}
	var parts []string
	for _, e := range v {
		parts = append(parts, e.Field+": "+e.Message)
	}
	return "validation failed: " + strings.Join(parts, "; ")
}
