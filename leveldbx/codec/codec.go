package codec

// LdbCodec leveldb codec interface
type LdbCodec[K any, V any] interface {
	// EncodeKey encode key
	EncodeKey(data K) (k []byte, err error)
	// DecodeKey decode key
	DecodeKey(data []byte) (k K, err error)
	// EncodeVal encode value
	EncodeVal(data V) (v []byte, err error)
	// DecodeVal decode value
	DecodeVal(data []byte) (v V, err error)
}
