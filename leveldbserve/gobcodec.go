package leveldbserve

import (
	"encoding/gob"

	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
)

// GobEncode gob encode
func GobEncode(data interface{}) (value []byte, err error) {
	defer rescueutil.Recover(func(err any) {
		err = err.(error)
	})
	buf := getEncodeBuf()
	defer putEncodeBuf(buf)
	enc := gob.NewEncoder(buf)
	errcheck.CheckAndPanic(enc.Encode(data))
	value = buf.Bytes()
	return
}

// GobDecode gob decod
func GobDecode(data []byte, to interface{}) (err error) {
	defer rescueutil.Recover(func(err any) {
		err = err.(error)
	})
	buf := getDecodeBuf()
	defer putDecodeBuf(buf)
	_ = assignutil.Assign(buf.Write(data))
	dec := gob.NewDecoder(buf)
	errcheck.CheckAndPanic(dec.Decode(to))
	return
}
