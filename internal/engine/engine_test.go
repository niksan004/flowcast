package engine

import (
	"reflect"
	"testing"
)

func TestRenderString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		env     map[string]any
		want    string
		wantErr bool
	}{
		{
			name:  "no placeholders",
			input: "hello render",
			env:   map[string]any{},
			want:  "hello render",
		},
		{
			name:  "one placeholder",
			input: "count is {{ count }}",
			env:   map[string]any{"count": 4},
			want:  "count is 4",
		},
		{
			name:  "one placeholder expression",
			input: "count is {{ count + 1 }}",
			env:   map[string]any{"count": 4},
			want:  "count is 5",
		},
		{
			name:  "two placeholders",
			input: "{{ object }} count is {{ count }}",
			env:   map[string]any{"object": "VM", "count": 4},
			want:  "VM count is 4",
		},
		{
			name:    "undefined variable",
			input:   "{{ missing }}",
			env:     map[string]any{},
			wantErr: true,
		},
		{
			name:  "resolve to int",
			input: "{{ count }}",
			env:   map[string]any{"count": 1},
			want:  "1",
		},
		{
			name:    "bad syntax",
			input:   "{{ count + }}",
			env:     map[string]any{"count": 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderVal(tt.input, tt.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderData(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		env     map[string]any
		want    map[string]any
		wantErr bool
	}{
		{
			name:  "render string field",
			input: map[string]any{"cmd": "echo hello {{ word }}"},
			env:   map[string]any{"word": "render"},
			want:  map[string]any{"cmd": "echo hello render"},
		},
		{
			name:  "skip int field",
			input: map[string]any{"value": 4},
			env:   map[string]any{},
			want:  map[string]any{"value": 4},
		},
		{
			name:  "skip field from list",
			input: map[string]any{"condition": "{{ count }} > 4"},
			env:   map[string]any{"count": 5},
			want:  map[string]any{"condition": "{{ count }} > 4"},
		},
		{
			name:    "failed to render string",
			input:   map[string]any{"cmd": "echo hello {{ word }}"},
			env:     map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderData(tt.input, tt.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
