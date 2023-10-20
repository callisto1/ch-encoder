package ch_encoding

// StringEncoder Encode/Decode Data via strings.
// Though it seems awfully unnecessary to force any implementation to then eventually use UTF-8 ch_encoding to both read/write. I think using
// ByteEncoder should be a better interface imo, to keep the data as bytes and avoid doubly ch_encoding.
type StringEncoder interface {
	Encode(data Data) (string, error)
	Decode(string) (Data, error)
}

// ByteArrayEncoder Encode and Decode Data via byte arrays.
type ByteArrayEncoder interface {
	EncodeBytes(data Data) ([]byte, error)
	DecodeBytes([]byte) (Data, error)
}
