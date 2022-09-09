package leveldbx

import "testing"

func TestOpenDB(t *testing.T) {
	db, err := OpenDB[string, string]("test")
	defer db.Close()
	if err != nil {
		t.Errorf("OpenDB() failed, err: %+v", err)
	}
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
