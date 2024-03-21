package batchprocess

import (
	"fmt"
	"testing"
)

func Test_batchProcess_Process(t *testing.T) {
	type fields[T interface{}, S interface{}] struct {
		items     []T
		batchSize int
		batchCnt  int
		result    []S
	}

	type test[T interface{}, S interface{}] struct {
		name     string
		fields   fields[T, S]
		function func(item T) S
	}
	var tests = []test[string, int]{
		{
			name: "test1",
			fields: fields[string, int]{
				items:    []string{"I", "am", "Iron", "man"},
				batchCnt: 2,
			},
			function: func(item string) int { return len(item) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBatchProcessor[string, int](tt.fields.items, tt.fields.batchCnt)
			p.Process(tt.function)
			for _, v := range p.result {
				fmt.Println(v)
			}
		})
	}
}
