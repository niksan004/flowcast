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

		// skip non-string fields
		stringVal, ok := val.(string)
		if !ok {
			renderedData[key] = val
			continue
		}

		// render string fields
		renderedString, err := renderString(stringVal, env)
		if err != nil {
			return nil, err
		}
		renderedData[key] = renderedString
	}

	return renderedData, nil
}

func renderString(s string, env map[string]any) (string, error) {
	var evalErr error

	out := regex.ReplaceAllStringFunc(s, func(match string) string {
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

	return out, evalErr
}
