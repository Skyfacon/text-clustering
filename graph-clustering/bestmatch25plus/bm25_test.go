package bestmatch25plus

import "testing"

func Test_dictFunc(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				words: []string{"I", "am", "Iron", "Man"},
			},
		},

		{
			name: "test2",
			args: args{
				words: []string{"I", "am", "Who", "I", "am"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dictFunc(tt.args.words)
		})
	}
}
