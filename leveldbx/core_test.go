package leveldbx

import (
	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec/gobc"
	"testing"
)

func TestDB_OpenDB(t *testing.T) {
	db, err := OpenDB[string, string]("test")
	defer func(db *DB[string, string]) {
		if err := db.Close(); err != nil {
			t.Errorf("dbClose error: %+v", err)
		}
		t.Logf("db: [%s] closed", "test")
	}(db)
	if err != nil {
		t.Errorf("OpenDB() failed, err: %+v", err)
	}
	k, _ := db.codec.EncodeKey("testKey1")
	t.Logf("putKey: %s", string(k))
	err = db.Put("testKey1", "abc")
	if err != nil {
		t.Errorf("Put() failed, err: %+v", err)
	}
	val, err := db.Get("testKey1")
	if err != nil {
		t.Errorf("Get() failed, err: %+v", err)
	} else {
		t.Logf("Get() success, val: %s", val)
	}
}

func TestDB_SetCodec(t *testing.T) {
	db, err := OpenDB[string, string]("test")
	defer func(db *DB[string, string]) {
		if err := db.Close(); err != nil {
			t.Errorf("dbClose error: %+v", err)
		}
		t.Logf("db: [%s] closed", "test")
	}(db)
	if err != nil {
		t.Errorf("OpenDB() failed, err: %+v", err)
	}
	db.SetCodec(gobc.NewGobCodec[string, string]())
	err = db.Put("testKey1", "abc")
	if err != nil {
		t.Errorf("Put() failed, err: %+v", err)
	}
	val, err := db.Get("testKey1")
	if err != nil {
		t.Errorf("Get() failed, err: %+v", err)
	} else {
		t.Logf("Get() success, val: %s", val)
	}
}
