package steps

import (
	"testing"
)

func TestIfStep(t *testing.T) {
	tests := []struct {
		name      string
		inputCond string
		inputThen []RawStep
		inputElse []RawStep
		env       map[string]any
		want      []RawStep
		wantErr   bool
	}{
		{
			name:      "true condition",
			inputCond: "count > 4",
			inputThen: []RawStep{RawStep{Name: "then"}},
			inputElse: []RawStep{RawStep{Name: "else"}},
			env:       map[string]any{"count": 5},
			want:      nil,
		},
		{
			name:      "false condition",
			inputCond: "count > 4",
			inputThen: []RawStep{RawStep{Name: "then"}},
			inputElse: []RawStep{RawStep{Name: "else"}},
			env:       map[string]any{"count": 3},
			want:      nil,
		},
		{
			name:      "non-bool condition",
			inputCond: "count",
			inputThen: []RawStep{RawStep{Name: "then"}},
			inputElse: []RawStep{RawStep{Name: "else"}},
			env:       map[string]any{"count": 3},
			wantErr:   true,
		},
		{
			name:      "bad syntax",
			inputCond: "count < ",
			inputThen: []RawStep{RawStep{Name: "then"}},
			inputElse: []RawStep{RawStep{Name: "else"}},
			env:       map[string]any{"count": 3},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := IfStep{Condition: tt.inputCond, Then: tt.inputThen, Else: tt.inputElse}
			err := step.Run(tt.env, func(nextSteps []RawStep) error {
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestForStep(t *testing.T) {
	tests := []struct {
		name      string
		inputCond string
		inputDo   []RawStep
		env       map[string]any
		want      []RawStep
		wantErr   bool
	}{
		{
			name:      "true condition",
			inputCond: "count > 4",
			inputDo:   []RawStep{RawStep{Name: "do"}},
			env:       map[string]any{"count": 5},
			want:      nil,
		},
		{
			name:      "false condition",
			inputCond: "count > 4",
			inputDo:   []RawStep{RawStep{Name: "do"}},
			env:       map[string]any{"count": 3},
			want:      nil,
		},
		{
			name:      "non-bool condition",
			inputCond: "count",
			inputDo:   []RawStep{RawStep{Name: "do"}},
			env:       map[string]any{"count": 3},
			wantErr:   true,
		},
		{
			name:      "bad syntax",
			inputCond: "count < ",
			inputDo:   []RawStep{RawStep{Name: "do"}},
			env:       map[string]any{"count": 3},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := ForStep{Condition: tt.inputCond, Do: tt.inputDo}
			err := step.Run(tt.env, func(nextSteps []RawStep) error {
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
