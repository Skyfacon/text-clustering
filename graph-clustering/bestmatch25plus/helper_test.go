package bestmatch25plus

import "testing"

func Test_wordSegInGo(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test01",
			args: args{s: "我有一个梦想，那就是成为海贼王"},
		},
		{
			name: "test02",
			args: args{s: "比特币这波牛市，大家一定要抓住机会"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wordSegInGo(tt.args.s)
		})
	}
}

func Test_simulate(t *testing.T) {
	type args struct {
		sentences []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{sentences: []string{"我是中国人，我爱中国", "我是炎黄子孙，我爱中国人民", "海贼王路飞吃了尼卡果实"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulate(tt.args.sentences)
		})
	}
}
