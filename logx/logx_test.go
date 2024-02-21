package logx

import "testing"

func TestNewLoggerOption(t *testing.T) {
	type args struct {
		opt []OptionFunc
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewLoggerOption(tt.args.opt...)
		})
	}
	Info("Hello")
}
