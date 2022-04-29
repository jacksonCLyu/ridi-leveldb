package leveldbserve

import "testing"

func TestGetDB(t *testing.T) {
	db := GetDB("test")
	if db == nil {
		t.Error("GetDB() failed")
	}
}
