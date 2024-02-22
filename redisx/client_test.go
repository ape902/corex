package redisx

import (
	"context"
	"testing"
)

func TestInitRedis(t *testing.T) {
	type args struct {
		addr   string
		pass   string
		db     int
		option []OptionFunc
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"Redis",
			args{
				addr: "127.0.0.1:6379",
				pass: "123456",
				db:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRedis(tt.args.addr, tt.args.pass, tt.args.db, tt.args.option...)

			if err := Client.Ping(context.Background()).Err(); err != nil {
				t.Error(err)
			}
		})
	}

}
