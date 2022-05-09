package leveldbserve

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/convutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var dbMemo *dbCache

var once sync.Once

// init logger & start db monitor
func init() {
	once.Do(func() {
		dbMemo = &dbCache{DbReqs: make(chan *dbRequest)}
		go dbMemo.serve(createDB)
	})
}

// Close To close the opened leveldb
func (db *UnderlineDB) Close() error {
	return db.background.Close()
}

// createDB create a db instance realtime
func createDB(path string, dbName string, options *opt.Options) (db *UnderlineDB, err error) {
	defer rescueutil.Recover(func(e any) {
		err = e.(error)
	})
	dbFileName := filepath.Join(path, DbRelativePath, dbName)
	ldb := assignutil.Assign(leveldb.OpenFile(dbFileName, options))
	db = &UnderlineDB{background: ldb}
	return
}

// Get get value from db by key
func (db *UnderlineDB) Get(key string) (value []byte) {
	defer rescueutil.Recover(func(err any) {
		fmt.Printf("Get Error: %v\n", err)
		value = nil
	})
	value = assignutil.Assign(db.background.Get(convutil.Str2bytes(key), &opt.ReadOptions{}))
	return
}

// Put put (k, v) into db
func (db *UnderlineDB) Put(key string, val interface{}) (err error) {
	defer rescueutil.Recover(func(e any) {
		err = e.(error)
	})
	v := assignutil.Assign(GobEncode(val))
	errcheck.CheckAndPanic(db.background.Put(convutil.Str2bytes(key), v, &opt.WriteOptions{}))
	return
}

// PutBatch batch put to db
func (db *UnderlineDB) PutBatch(rsMap map[string]interface{}) (err error) {
	defer rescueutil.Recover(func(e any) {
		err = e.(error)
	})
	batch := new(leveldb.Batch)
	for k, v := range rsMap {
		value, err := GobEncode(v)
		if err != nil {
			continue
		}
		batch.Put(convutil.Str2bytes(k), value)
	}
	return db.background.Write(batch, nil)
}

// DeleteBatch batch delete
func (db *UnderlineDB) DeleteBatch(keys ...string) (err error) {
	defer rescueutil.Recover(func(e any) {
		err = e.(error)
	})
	if len(keys) == 0 {
		return
	}
	batch := new(leveldb.Batch)
	for _, key := range keys {
		batch.Delete(convutil.Str2bytes(key))
	}
	return db.background.Write(batch, nil)
}

// ForRange traverse
func (db *UnderlineDB) ForRange(applyFunc func(k string, v []byte)) {
	defer rescueutil.Recover(func(e any) {
		fmt.Printf("ForRange Error: %v\n", e)
	})
	if nil == applyFunc {
		return
	}
	it := db.background.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		key := convutil.Bytes2str(it.Key())
		value := it.Value()
		applyFunc(key, value)
	}
}

// ForRangeAsync traverse
func (db *UnderlineDB) ForRangeAsync(applyFunc func(k string, v []byte)) {
	defer rescueutil.Recover(func(e any) {
		fmt.Printf("ForRangeAsync Error: %v\n", e)
	})
	if nil == applyFunc {
		return
	}
	ch := make(chan *DbEntry)
	go db.forRangeChan(ch)

	for e := range ch {
		applyFunc(e.Key, e.Value)
	}
}

func (db *UnderlineDB) forRangeChan(enChan chan *DbEntry) {
	defer rescueutil.Recover(func(e any) {
		fmt.Printf("ForRangeChan Error: %v\n", e)
	})
	it := db.background.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		entry := &DbEntry{
			Key:   convutil.Bytes2str(it.Key()),
			Value: it.Value(),
		}
		enChan <- entry
	}
	close(enChan)
}
