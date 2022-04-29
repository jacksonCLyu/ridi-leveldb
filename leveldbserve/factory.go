package leveldbserve

import (
	"github.com/jacksonCLyu/ridi-faces/pkg/env"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// GetDB The entrance of get db request by default home dir
func GetDB(dbName string) *underlineDB {
	return GetDBWithPathAndOpts(env.AppRootPath(), dbName, &opt.Options{})
}

// GetDBWithOpts The entrance of get db request by default home dir and *opt.Options
func GetDBWithOpts(dbName string, opts *opt.Options) *underlineDB {
	return GetDBWithPathAndOpts(env.AppRootPath(), dbName, opts)
}

// GetDBWithPathAndOpts The entrance of get db request
func GetDBWithPathAndOpts(path string, dbName string, opts *opt.Options) *underlineDB {
	if "" == dbName {
		return nil
	}
	dbReq := &dbRequest{
		Path:   path,
		DbName: dbName,
		Opts:   opts,
		DbChan: make(chan *underlineDB),
	}
	dbMemo.DbReqs <- dbReq
	return <-dbReq.DbChan
}
