package leveldbserve

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/jacksonCLyu/ridi-faces/pkg/env"
	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"path/filepath"
)

type DB[K any, V any] struct {
	underlineDB   *leveldb.DB
	localARCCache *lru.ARCCache
}

func NewDB[K any, V any](name string, size int, options *opt.Options) (*DB[K, V], error) {
	dbFileName := filepath.Join(env.AppRootPath(), DbRelativePath, name)
	ldb := assignutil.Assign(leveldb.OpenFile(dbFileName, options))
	cache, err := lru.NewARC(size)
	if err != nil {
		return nil, err
	}
	return &DB[K, V]{
		underlineDB:   ldb,
		localARCCache: cache,
	}, nil
}
