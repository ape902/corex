package redisx

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var Client *redis.Client

type (
	options struct {
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		PoolSize     int
		PoolTimeout  time.Duration
		IdleTimeout  time.Duration
	}

	OptionFunc func(opt *options)
)

func WithDialTimeout(t time.Duration) OptionFunc {
	return func(opt *options) {
		opt.DialTimeout = t
	}
}
func WithReadTimeout(t time.Duration) OptionFunc {
	return func(opt *options) {
		opt.ReadTimeout = t
	}
}
func WithWriteTimeout(t time.Duration) OptionFunc {
	return func(opt *options) {
		opt.WriteTimeout = t
	}
}
func WithPoolSize(size int) OptionFunc {
	return func(opt *options) {
		opt.PoolSize = size
	}
}
func WithPoolTimeout(t time.Duration) OptionFunc {
	return func(opt *options) {
		opt.PoolTimeout = t
	}
}
func WithIdleTimeout(t time.Duration) OptionFunc {
	return func(opt *options) {
		opt.IdleTimeout = t
	}
}

func InitRedis(addr, pass string, db int, option ...OptionFunc) {
	opts := &options{}
	for _, o := range option {
		if o != nil {
			o(opts)
		}
	}

	opt := &redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           db,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolSize:     opts.PoolSize,
		PoolTimeout:  opts.PoolTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}

	Client = redis.NewClient(opt)
}
