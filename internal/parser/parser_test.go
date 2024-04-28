package parser_test

import (
	"testing"

	"github.com/zeet-co/kang/internal/parser"
)

func TestGetValue(t *testing.T) {
	type Nested struct {
		Detail string `json:"detail"`
	}
	type TestStruct struct {
		Name    string            `json:"name"`
		Numbers []int             `json:"numbers"`
		Nested  Nested            `json:"nested"`
		Info    map[string]string `json:"info"`
	}

	obj := TestStruct{
		Name:    "Test",
		Numbers: []int{1, 2, 3},
		Nested:  Nested{Detail: "Detailed Info"},
		Info:    map[string]string{"key1": "value1", "key2": "value2"},
	}

	tests := []struct {
		name     string
		jsonPath string
		want     string
	}{
		{
			name:     "Read simple field",
			jsonPath: "name",
			want:     "Test",
		},
		{
			name:     "Read nested field",
			jsonPath: "nested.detail",
			want:     "Detailed Info",
		},
		{
			name:     "Read array field",
			jsonPath: "numbers[0]",
			want:     "1",
		},
		{
			name:     "Read array field",
			jsonPath: "numbers[1]",
			want:     "2",
		},
		{
			name:     "Read map field",
			jsonPath: "info.key1",
			want:     "value1",
		},
		{
			name:     "Read undefined field",
			jsonPath: "undefined",
			want:     "",
		},
		{
			name:     "Read out-of-bounds array index",
			jsonPath: "numbers[5]",
			want:     "", // or the expected error message or behavior
		},
		{
			name:     "Read non-existent map key",
			jsonPath: "info.key3",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.GetValue(&obj, tt.jsonPath); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
