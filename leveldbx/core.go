package leveldbx

import (
	"errors"
	"path/filepath"

	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec"
	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec/jsonc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DB[K any, V any] struct {
	underlineDB *leveldb.DB
	codec       codec.LCodec[K, V]
}

// OpenDB open a leveldb with the given name and opts
// As usual, when you call OpenDB(), `defer db.close ()` is used to Close it
//
// To use bloom filter:
//
//	leveldbx.OpenDB("foo", leveldbx.WithLevelDBOpts(&opt.Options{Filter: filter.NewBloomFilter(10),}))
func OpenDB[K any, V any](name string, opts ...Option) (*DB[K, V], error) {
	options := DefaultOptions()
	for _, leveldbOpt := range opts {
		leveldbOpt.apply(options)
	}
	file := filepath.Join(options.path, name)
	ldb, err := leveldb.OpenFile(file, options.ldbOpts)
	if err != nil {
		return nil, err
	}
	return &DB[K, V]{
		underlineDB: ldb,
		codec:       jsonc.NewJSONCodec[K, V](),
	}, nil
}

// Close to close level db
func (db *DB[K, V]) Close() error {
	return db.underlineDB.Close()
}

// SetCodec set leveldb codec
func (db *DB[K, V]) SetCodec(codec codec.LCodec[K, V]) {
	db.codec = codec
}

// Get key from leveldb
func (db *DB[K, V]) Get(key K) (v V, err error) {
	var k []byte
	k, err = db.codec.EncodeKey(key)
	if err != nil {
		return
	}
	var value []byte
	value, err = db.underlineDB.Get(k, &opt.ReadOptions{})
	if err != nil {
		return
	}
	v, err = db.codec.DecodeVal(value)
	return
}

// GetWithReadOpts get key from leveldb with read options
func (db *DB[K, V]) GetWithReadOpts(key K, readOpts *opt.ReadOptions) (v V, err error) {
	var k []byte
	k, err = db.codec.EncodeKey(key)
	if err != nil {
		return
	}
	var value []byte
	value, err = db.underlineDB.Get(k, readOpts)
	if err != nil {
		return
	}
	v, err = db.codec.DecodeVal(value)
	return
}

// Put key value pair to leveldb
func (db *DB[K, V]) Put(key K, val V) error {
	k, err := db.codec.EncodeKey(key)
	if err != nil {
		return err
	}
	v, err := db.codec.EncodeVal(val)
	if err != nil {
		return err
	}
	return db.underlineDB.Put(k, v, &opt.WriteOptions{})
}

// PutWithWriteOpts put key value pair to leveldb with write options
func (db *DB[K, V]) PutWithWriteOpts(key K, val V, wo *opt.WriteOptions) error {
	k, err := db.codec.EncodeKey(key)
	if err != nil {
		return err
	}
	v, err := db.codec.EncodeVal(val)
	if err != nil {
		return err
	}
	return db.underlineDB.Put(k, v, wo)
}

func (db *DB[K, V]) EncodeKey(k K) ([]byte, error) {
	return db.codec.EncodeKey(k)
}

func (db *DB[K, V]) EncodeVal(v V) ([]byte, error) {
	return db.codec.EncodeVal(v)
}

func (db *DB[K, V]) DecodeKey(data []byte) (K, error) {
	return db.codec.DecodeKey(data)
}

func (db *DB[K, V]) DecodeVal(data []byte) (V, error) {
	return db.codec.DecodeVal(data)
}

func (db *DB[K, V]) PutBatch(batch *leveldb.Batch) error {
	return db.underlineDB.Write(batch, &opt.WriteOptions{})
}

func (db *DB[K, V]) PutBatchWithWriteOpts(batch *leveldb.Batch, wo *opt.WriteOptions) error {
	return db.underlineDB.Write(batch, wo)
}

func (db *DB[K, V]) DeleteBatch(keys ...K) error {
	if len(keys) == 0 {
		return nil
	}
	errKeys := make([]K, 0)
	batch := new(leveldb.Batch)
	for _, key := range keys {
		k, err := db.codec.EncodeKey(key)
		if err != nil {
			errKeys = append(errKeys, key)
		}
		batch.Delete(k)
	}
	return db.underlineDB.Write(batch, &opt.WriteOptions{})
}

func (db *DB[K, V]) DeleteBatchWithWriteOpts(keys []K, wo *opt.WriteOptions) error {
	if len(keys) == 0 {
		return nil
	}
	errKeys := make([]K, 0)
	batch := new(leveldb.Batch)
	for _, key := range keys {
		k, err := db.codec.EncodeKey(key)
		if err != nil {
			errKeys = append(errKeys, key)
		}
		batch.Delete(k)
	}
	return db.underlineDB.Write(batch, wo)
}

// ForRange add for range
func (db *DB[K, V]) ForRange(consumer func(key K, val V)) error {
	return db.ForRangeWithOpts(nil, nil, consumer)
}

// Seek seek startKey, then for range
func (db *DB[K, V]) Seek(startKey K, consumer func(key K, val V)) error {
	if consumer == nil {
		return errors.New("ForRange consumer can not be nil")
	}
	sk, err := db.codec.EncodeKey(startKey)
	if err != nil {
		return err
	}
	it := db.underlineDB.NewIterator(nil, nil)
	for ok := it.Seek(sk); ok; ok = it.Next() {
		// 忽略解码错误，完成遍历语义
		k, _ := db.codec.DecodeKey(it.Key())
		v, _ := db.codec.DecodeVal(it.Value())
		consumer(k, v)
	}
	it.Release()
	return it.Error()
}

// ForRangeWithOpts for range with *util.Range and options
// *util.Range can be:
//
//	&util.Range{Start: []byte("foo"), Limit: []byte("xoo")}
//
// or can be:
//
//	util.BytesPrefix([]byte("foo-")
func (db *DB[K, V]) ForRangeWithOpts(slice *util.Range, ro *opt.ReadOptions, consumer func(key K, val V)) error {
	if consumer == nil {
		return errors.New("ForRange consumer can not be nil")
	}
	it := db.underlineDB.NewIterator(slice, ro)
	for it.Next() {
		// 忽略解码错误，完成遍历语义
		k, _ := db.codec.DecodeKey(it.Key())
		v, _ := db.codec.DecodeVal(it.Value())
		consumer(k, v)
	}
	it.Release()
	return it.Error()
}

// Size returns db size
func (db *DB[K, V]) Size() int {
	size := 0
	if err := db.ForRange(func(_ K, _ V) {
		size++
	}); err != nil {
		return size
	}
	return size
}
