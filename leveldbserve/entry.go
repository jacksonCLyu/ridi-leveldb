package leveldbserve

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// UnderlineDB struct of leveldb.UnderlineDB
type UnderlineDB struct {
	background *leveldb.DB
}

// Create a new local db if none exists.
// Otherwise, open the existing db
type newDB func(path string, dbName string, options *opt.Options) (*UnderlineDB, error)

// cacheResult cache result of local db
type cacheResult struct {
	underlineDB *UnderlineDB
	ready       chan struct{}
}

// dbRequest The async request of get db
type dbRequest struct {
	Path   string
	DbName string
	Opts   *opt.Options
	DbChan chan *UnderlineDB
}

// dbCache leveldb cache
type dbCache struct {
	DbReqs chan *dbRequest
}

// DbEntry The entry of leveldb key value pair
type DbEntry struct {
	Key   string
	Value []byte
}
