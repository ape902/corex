package gormx

import (
	"testing"
)

func TestInitGorm(t *testing.T) {
	type args struct {
		dbType  string
		user    string
		pass    string
		host    string
		port    string
		dbName  string
		optFunc []OptionFunc
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"gorm",
			args{
				dbType: "mysql",
				user:   "root",
				pass:   "123456",
				host:   "127.0.0.1",
				port:   "3306",
				dbName: "test",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitGorm(tt.args.dbType, tt.args.user, tt.args.pass, tt.args.host, tt.args.port, tt.args.dbName, tt.args.optFunc...)
		})
	}
}
