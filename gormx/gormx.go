package gormx

import (
	"errors"
	"fmt"
	"github.com/ape902/corex/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

const (
	DefaultMaxOpenConn       = 1000
	DefaultMaxIdleConn       = 100
	DefaultConnMaxLifeSecond = 30 * time.Minute
	DefaultLogName           = "gorm"
)

type (
	gormStruct struct {
		dbType   string // 数据库驱动类型（Default: mysql）
		username string // 用户名
		password string // 密码
		host     string // 主机
		port     string // 端口
		dbName   string // 数据库名
	}

	option struct {
		MaxOpenConn       int
		MaxIdleConn       int
		ConnMaxLifeSecond time.Duration
		PrepareStmt       bool
		LogName           string
	}

	OptionFunc func(opt *option)
)

func (opt *option) WithMaxOpenConn(num int) OptionFunc {
	return func(opt *option) {
		opt.MaxOpenConn = num
	}
}

func (opt *option) WithMaxIdleConn(num int) OptionFunc {
	return func(opt *option) {
		opt.MaxIdleConn = num
	}
}

func (opt *option) WithMaxLifeSecond(num time.Duration) OptionFunc {
	return func(opt *option) {
		opt.ConnMaxLifeSecond = num
	}
}

func (opt *option) WithLogName(name string) OptionFunc {
	return func(opt *option) {
		opt.LogName = name
	}
}

var Client *gorm.DB

func InitGorm(dbType, user, pass, host, port, dbName string, optFunc ...OptionFunc) {
	logx.NewLoggerOption()
	opt := &option{}
	for _, f := range optFunc {
		if f != nil {
			f(opt)
		}
	}
	if opt.ConnMaxLifeSecond == 0 {
		opt.ConnMaxLifeSecond = DefaultConnMaxLifeSecond
	}
	if opt.MaxIdleConn == 0 {
		opt.MaxIdleConn = DefaultMaxIdleConn
	}

	var gs gormStruct
	gs.dbType = dbType
	gs.username = user
	gs.password = pass
	gs.host = host
	gs.port = port
	gs.dbName = dbName

	dsn := gs.BuildDSN()
	db, err := gs.gormDial(dsn, opt)
	if err != nil {
		logx.Panic(err)
	}

	Client = db
}

func (g *gormStruct) BuildDSN() gorm.Dialector {
	switch g.dbType {
	case "mysql":
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
			g.username,
			g.password,
			g.host,
			g.port,
			g.dbName,
			true,
			"Local"))
	default:
		return nil
	}

}

func (g *gormStruct) gormDial(dial gorm.Dialector, option *option) (*gorm.DB, error) {
	db, err := gorm.Open(dial, &gorm.Config{
		//为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）
		//如果没有这方面的要求，可以设置SkipDefaultTransaction为true来禁用它。
		//SkipDefaultTransaction: true,
		//Logger: Log,
		//执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续执行的效率
		PrepareStmt: option.PrepareStmt,
		NamingStrategy: schema.NamingStrategy{
			//使用单数表名,默认为复数表名，即当model的结构体为User时，默认操作的表名为users
			//设置	SingularTable: true 后当model的结构体为User时，操作的表名为user
			SingularTable: true,

			//TablePrefix: "pre_", //表前缀
		},
		//Logger: logx.Default.LogMode(logx.Info), // 日志配置
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s, [db connection failed] Database name: %s", err, g.dbName))
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	if option.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(option.MaxOpenConn)
	} else {
		sqlDB.SetMaxOpenConns(DefaultMaxOpenConn)
	}

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	if option.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(option.MaxIdleConn)
	}

	// 设置最大连接超时时间
	if option.ConnMaxLifeSecond > 0 {
		sqlDB.SetConnMaxLifetime(time.Second * option.ConnMaxLifeSecond)
	}

	err = db.Callback().Create().After("gorm:after_create").Register(DefaultLogName, afterLog)
	if err != nil {
		logx.Errorf("Register Create error, %s", err)
	}
	err = db.Callback().Query().After("gorm:after_query").Register(DefaultLogName, afterLog)
	if err != nil {
		logx.Errorf("Register Query error", err)
	}
	err = db.Callback().Update().After("gorm:after_update").Register(DefaultLogName, afterLog)
	if err != nil {
		logx.Errorf("Register Update error", err)
	}
	err = db.Callback().Delete().After("gorm:after_delete").Register(DefaultLogName, afterLog)
	if err != nil {
		logx.Errorf("Register Delete error", err)
	}
	return db, nil
}

func afterLog(db *gorm.DB) {
	err := db.Error
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	if err != nil {
		logx.Error(sql, err)
	} else {
		logx.Infof("[ SQL语句: %s]", sql)
	}
}
