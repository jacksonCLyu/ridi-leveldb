package leveldbserve

import (
	"fmt"
	"sync"
	"time"

	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var _dbCacheMap sync.Map

func getDBCacheOrDefault(dbName string, cache *cacheResult) (*cacheResult, bool) {
	db, loaded := _dbCacheMap.LoadOrStore(dbName, cache)
	return db.(*cacheResult), loaded
}

func (*dbCache) serve(f newDB) {
	//fmt.Printf("%s\t%s\tDbMemo Monitor Serve Goroutine Start.\n", time.Now().Local().Format(env.ISO8601UTC), "INFO")
	for dbReq := range dbMemo.DbReqs {
		r, loaded := getDBCacheOrDefault(dbReq.DbName, &cacheResult{
			ready: make(chan struct{}),
		})
		if !loaded {
			// No cache found, call create a new one
			go r.call(f, dbReq.Path, dbReq.DbName, dbReq.Opts)
		} else {
			// Cache found, check if Closed then reopen db
			s, e := r.underlineDB.background.GetSnapshot()
			if s == nil || e == leveldb.ErrClosed {
				r = &cacheResult{
					ready: make(chan struct{}),
				}
				go r.call(f, dbReq.Path, dbReq.DbName, dbReq.Opts)
			}
		}
		// As long as the call is initiated, there is no need to wait for the next call, whether it is created or reopened
		go r.deliver(dbReq.DbChan)
	}
}

func (r *cacheResult) call(f newDB, path, dbName string, options *opt.Options) {
	defer rescueutil.Recover(func(err any) {
		if e, ok := err.(error); ok {
			fmt.Printf("%s\t%s\tDbMemo Monitor Call Goroutine Error: %s\n", time.Now().Local().Format(time.RFC3339), "ERROR", e.Error())
		} else {
			fmt.Printf("%s\t%s\tDbMemo Monitor Call Goroutine Error: %v\n", time.Now().Local().Format(time.RFC3339), "ERROR", err)
		}
	})
	r.underlineDB = assignutil.Assign(f(path, dbName, options))
	close(r.ready)
}

func (r *cacheResult) deliver(dbChan chan *UnderlineDB) {
	<-r.ready
	dbChan <- r.underlineDB
}
