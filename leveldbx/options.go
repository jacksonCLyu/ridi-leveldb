package leveldbx

import (
	"path/filepath"

	"github.com/jacksonCLyu/ridi-faces/pkg/env"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Option interface {
	apply(opts *LdbOptions)
}

type ApplyFunc func(opts *LdbOptions)

func (f ApplyFunc) apply(opts *LdbOptions) {
	f(opts)
}

// LdbOptions private LdbOptions
type LdbOptions struct {
	// path dir path of db files
	path string
	// ldbOpts leveldb LdbOptions
	ldbOpts *opt.Options
}

func DefaultOptions() *LdbOptions {
	return &LdbOptions{
		path:    filepath.Join(env.AppRootPath(), DbRelativePath),
		ldbOpts: &opt.Options{},
	}
}

// WithPath set dir path of db files
func WithPath(path string) Option {
	return ApplyFunc(func(opts *LdbOptions) {
		opts.path = path
	})
}

// WithLevelDBOpts set leveldb opt.Options
func WithLevelDBOpts(levelDBOpts *opt.Options) Option {
	return ApplyFunc(func(opts *LdbOptions) {
		opts.ldbOpts = levelDBOpts
	})
}
