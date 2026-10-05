package engine

import (
	"fmt"
	"github.com/expr-lang/expr"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`\{\{.*?\}\}`)

var skipFields = map[string]bool{
	"condition": true,
	"do":        true,
	"then":      true,
	"else":      true,
}

func renderData(data map[string]any, env map[string]any) (map[string]any, error) {
	renderedData := make(map[string]any, len(data))
	for key, val := range data {
		// skip fields we do not want to render
		if skipFields[key] {
			renderedData[key] = val
			continue
		}

		// render fields
		renderedString, err := renderVal(val, env)
		if err != nil {
			return nil, err
		}
		renderedData[key] = renderedString
	}

	return renderedData, nil
}

func renderVal(val any, env map[string]any) (any, error) {
	switch val := val.(type) {
	case string:
		var evalErr error
		out := regex.ReplaceAllStringFunc(val, func(match string) string {
			expression := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}"))
			val, err := expr.Eval(expression, env)
			// report first error only
			if evalErr != nil {
				return ""
			}
			if err != nil {
				evalErr = err
				return ""
			}
			if val == nil {
				evalErr = fmt.Errorf("Undefined variable while evaluating %s", expression)
				return ""
			}
			return fmt.Sprint(val)
		})
		return out, nil
	// recurse into maps
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v := range val {
			renderedV, err := renderVal(v, env)
			if err != nil {
				return nil, err
			}
			out[k] = renderedV
		}
		return out, nil
	// recurse into slices
	case []any:
		out := make([]any, len(val))
		for i, v := range val {
			renderedV, err := renderVal(v, env)
			if err != nil {
				return nil, err
			}
			out[i] = renderedV
		}
		return out, nil
	default:
		return val, nil
	}
}
