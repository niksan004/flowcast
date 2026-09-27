package steps

import (
	"reflect"
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
			want:      []RawStep{RawStep{Name: "then"}},
		},
		{
			name:      "false condition",
			inputCond: "count > 4",
			inputThen: []RawStep{RawStep{Name: "then"}},
			inputElse: []RawStep{RawStep{Name: "else"}},
			env:       map[string]any{"count": 3},
			want:      []RawStep{RawStep{Name: "else"}},
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
			got, err := step.Branch(tt.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
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
			want:      []RawStep{RawStep{Name: "do"}},
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
			got, err := step.Loop(tt.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
