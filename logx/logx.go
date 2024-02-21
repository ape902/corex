package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	log      *zap.SugaredLogger
	levelMap = map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
	setupOnce sync.Once
	fileName  string
)

type (
	loggerConf struct {
		// Mode represents the logging mode, default is `console`.
		// console: log to console.
		// file: log to file.
		Mode string `json:",default=console,options=[console,file]"`
		// 日志路径
		LogPath string `json:",default=./logs"`
		// 日志级别
		Level string `json:",default=info,options=[Debug，Info，Warn，Error，DPanic，Panic，Fatal]"`
		// 日志文件名
		ServerName string `json:",default=zsops"`
		//日志文件大小，默认不限制。【单位MB】
		MaxSize int `json:",default=0"`
		//最长保存天数，默认不限制。
		MaxDay int `json:",default=0"`
		//最多备份几个，默认不限制。
		MaxBackups int `json:",default=0"`
		//是否压缩文件，默认关闭。【使用gzip】
		Compress bool
	}

	OptionFunc func(conf *loggerConf)
)

/*
WithMode represents the logging mode, default is `console`.

	console: log to console.
	file: log to file.
*/
func WithMode(mode string) OptionFunc {
	return func(conf *loggerConf) {
		conf.Mode = mode
	}
}

/*
WithLogPath 日志路径。默认目录：logs
*/
func WithLogPath(p string) OptionFunc {
	return func(conf *loggerConf) {
		conf.LogPath = p
	}
}

/*
WithLevel 日志级别。默认为info

	参数：Debug，Info，Warn，Error，DPanic，Panic，Fatal
*/
func WithLevel(level string) OptionFunc {
	strings.ToLower(level)
	return func(conf *loggerConf) {
		conf.Level = strings.ToLower(level)
	}
}

/*
WithServerName 业务服务名。默认：oneops
*/
func WithServerName(name string) OptionFunc {
	return func(conf *loggerConf) {
		conf.ServerName = name
	}
}

/*
WithMaxSize 日志文件大小。默认：0

	0：不限制
	单位：【单位MB】
*/
func WithMaxSize(size int) OptionFunc {
	return func(conf *loggerConf) {
		conf.MaxSize = size
	}
}

/*
WithMaxDay 最长保存天数。默认：0

	0：不限制
*/
func WithMaxDay(size int) OptionFunc {
	return func(conf *loggerConf) {
		conf.MaxDay = size
	}
}

/*
WithBackups 最多备份几个，默认：0

	0：不限制
*/
func WithBackups(size int) OptionFunc {
	return func(conf *loggerConf) {
		conf.MaxBackups = size
	}
}

/*
WithCompress 是否压缩文件。默认：false

	false: 停用
	true：启用
	压缩类型：gzip
*/
func WithCompress(status bool) OptionFunc {
	return func(conf *loggerConf) {
		conf.Compress = status
	}
}

func NewLoggerOption(opt ...OptionFunc) {
	conf := &loggerConf{}
	for _, f := range opt {
		if f != nil {
			f(conf)
		}
	}

	fileName = path.Join(conf.LogPath, conf.ServerName+".log")

	setupOnce.Do(func() {
		syncWriters := make([]zapcore.WriteSyncer, 0)
		switch conf.Mode {
		case "file":
			defer func() {
				if r := recover(); r != nil {
					syncWriters = append(syncWriters, zapcore.AddSync(os.Stdout))
				}
			}()

			if conf.LogPath == "" {
				conf.LogPath = "./log"
			}
			if conf.ServerName == "" {
				conf.ServerName = "log"
			}
			fileName = path.Join(conf.LogPath, conf.ServerName+".log")

			newLogger := &lumberjack.Logger{
				Filename:   fileName,        // 日志文件名
				MaxSize:    conf.MaxSize,    // 日志文件大小
				MaxAge:     conf.MaxDay,     // 最长保存天数
				MaxBackups: conf.MaxBackups, // 最多备份几个
				Compress:   conf.Compress,   // 是否压缩文件，使用gzip
			}
			syncWriters = append(syncWriters, zapcore.AddSync(newLogger))
		default:
			syncWriters = append(syncWriters, zapcore.AddSync(os.Stdout))
		}

		level := getLoggerLevel(conf.Level)
		encoder := zap.NewProductionEncoderConfig()
		encoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
		}
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoder),
			zapcore.NewMultiWriteSyncer(syncWriters...),
			zap.NewAtomicLevelAt(level))
		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		log = logger.Sugar()
	})
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func DPanic(args ...interface{}) {
	log.DPanic(args...)
}

func DPanicf(format string, args ...interface{}) {
	log.DPanicf(format, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
