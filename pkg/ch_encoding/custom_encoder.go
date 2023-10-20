package ch_encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	IntCodepoint    = 0x1
	StringCodepoint = 0x2
	DataCodepoint   = 0x3
	Int32Size       = 4
)

// CustomEncoder encodes Data into a string using following schema:
//  1. map the possible item types (int, string, Data) to a codepoint byte:
//     0x1 -> int32,
//     0x2 -> string
//     0x3 -> Data
//  2. write the data to a byte buffer in the following manner
//     a) first write the corresponding codepoint
//     b) then write the number of bytes the data needs (unless int, in which case it always 4). For the case of Data, this is recursively calculated.
//     c) write the data
//
// A stringified version of the byte buffer is returned because that's what the assignment interface expected, though I feel it should just be a []byte.
type CustomEncoder struct{}

func (e *CustomEncoder) Encode(data Data) (string, error) {
	b, err := e.EncodeBytes(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *CustomEncoder) Decode(s string) (Data, error) {
	d, err := e.DecodeBytes([]byte(s))
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (e *CustomEncoder) EncodeBytes(data Data) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	for _, d := range data {
		switch t := d.(type) {
		// same memory configuration
		case int32:
			buf.WriteByte(IntCodepoint)
			// uint32 always takes 4 bytes, so we don't need to write the length
			binary.Write(buf, binary.BigEndian, t)
		case string:
			buf.WriteByte(StringCodepoint)
			p := []byte(t)
			// since strings have max length of 1M, and use 4 bytes for each character that's at most 4M, we can express this using 4 bytes as well. (thought actually we only need 3 bytes, or even 22 bits if we wanted to optimize even further)
			binary.Write(buf, binary.BigEndian, int32(len(p)))
			buf.Write(p)
		case Data:
			s := t.GetEncodedByteSize()
			buf.WriteByte(DataCodepoint)
			binary.Write(buf, binary.BigEndian, s)
			encoded, err := e.Encode(t)
			if err != nil {
				return nil, err
			}
			buf.Write([]byte(encoded))
		default:
			return nil, fmt.Errorf("cannot Encode type %v", t)
		}
	}
	return buf.Bytes(), nil
}

func (e *CustomEncoder) DecodeBytes(s []byte) (Data, error) {
	ret := make(Data, 0)
	if len(s) == 0 {
		return ret, nil
	}
	br := bytes.NewReader(s)
	for {
		readByte, err := br.ReadByte()
		if err == io.EOF {
			// reached the end
			break
		}
		if err != nil {
			return nil, err
		}
		if readByte == IntCodepoint {
			b2 := make([]byte, Int32Size)
			br.Read(b2)
			n := binary.BigEndian.Uint32(b2)

			ret = append(ret, int32(n))
		} else if readByte == StringCodepoint {
			// get number of bytes of the string, expressed as a 4 byte number (uint32)
			b2 := make([]byte, Int32Size)
			br.Read(b2)
			n := binary.BigEndian.Uint32(b2)

			// read all the string data
			stringData := make([]byte, n)
			br.Read(stringData)
			ret = append(ret, string(stringData))
		} else if readByte == DataCodepoint {
			b2 := make([]byte, Int32Size)
			br.Read(b2)
			n := binary.BigEndian.Uint32(b2)

			stringData := make([]byte, n)
			br.Read(stringData)
			childData, err := e.Decode(string(stringData))
			if err != nil {
				return nil, err
			}
			ret = append(ret, childData)
		} else {
			return nil, fmt.Errorf("invalid byte string supplied %v", s)
		}
	}

	return ret, nil
}
