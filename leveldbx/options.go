package leveldbx

import (
	"path/filepath"

	"github.com/jacksonCLyu/ridi-faces/pkg/env"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Option interface {
	apply(opts *options)
}

type ApplyFunc func(opts *options)

func (f ApplyFunc) apply(opts *options) {
	f(opts)
}

// options private options
type options struct {
	// path dir path of db files
	path string
	// ldbOpts leveldb options
	ldbOpts *opt.Options
}

func DefaultOptions() *options {
	return &options{
		path:    filepath.Join(env.AppRootPath(), DbRelativePath),
		ldbOpts: &opt.Options{},
	}
}

// WithPath set dir path of db files
func WithPath(path string) Option {
	return ApplyFunc(func(opts *options) {
		opts.path = path
	})
}

// WithLevelDBOpts set leveldb opt.Options
func WithLevelDBOpts(levelDBOpts *opt.Options) Option {
	return ApplyFunc(func(opts *options) {
		opts.ldbOpts = levelDBOpts
	})
}
