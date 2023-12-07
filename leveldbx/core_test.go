package leveldbx

import (
	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec/gobc"
	"github.com/stretchr/testify/assert"
	"strconv"
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
	assert.EqualValues(t, val, "abc")
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
	assert.EqualValues(t, val, "abc")
}

func TestDB_Size(t *testing.T) {
	name1 := "test1"
	db1, err := OpenDB[any, any](name1)
	defer func(db1 *DB[any, any]) {
		err := db1.Close()
		if err != nil {

		}
	}(db1)
	if err != nil {
		t.Errorf("OpenDB() failed, err: %+v", err)
	}
	_ = db1.Put("1", "1")

	name2 := "test2"
	db2, err := OpenDB[any, any](name2)
	defer func(db2 *DB[any, any]) {
		err := db2.Close()
		if err != nil {

		}
	}(db2)
	if err != nil {
		t.Errorf("OpenDB() failed, err: %+v", err)
	}
	_ = db2.Put("1", "1")
	_ = db2.Put("2", "2")

	name3 := "test3"
	db3, err := OpenDB[any, any](name3)
	defer func(db3 *DB[any, any]) {
		err := db3.Close()
		if err != nil {

		}
	}(db3)
	if err != nil {
		t.Errorf("OpenDB() failed, err: %+v", err)
	}
	for i := 0; i < 10000; i++ {
		_ = db3.Put(strconv.Itoa(i), "1")
	}
	type testCase[K any, V any] struct {
		name string
		db   *DB[any, any]
		want int
	}
	tests := []testCase[any, any]{
		{name: name1, db: db1, want: 1},
		{name: name2, db: db2, want: 2},
		{name: name3, db: db3, want: 10000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
