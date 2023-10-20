package ch_encoding

import (
	"strconv"
	"strings"
)

type Data []interface{}

func (d Data) GetEncodedByteSize() int32 {
	count := int32(0)
	for _, item := range d {
		switch t := item.(type) {
		case int32:
			// int32 takes 4 bytes for the value + 1 byte for the codepoint identifier
			count += 4 + 1
		case string:
			// length of the data, plus 4 for the count, plus 1 for the codepoint identifier
			count += int32(len([]byte(t))) + 4 + 1
		case Data:
			// length of the data, plus 4 for the count, plus 1 for the codepoint identifier
			count += 1 + 4 + t.GetEncodedByteSize()
		}
	}
	return count
}

func (d Data) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, item := range d {
		if i != 0 {
			sb.WriteString(", ")
		}
		switch t := item.(type) {
		// same memory configuration
		case int32:
			sb.WriteString(strconv.Itoa(int(t)))
		case string:
			sb.WriteString(t)
		case Data:
			s := t.String()
			sb.WriteString(s)
		}
	}
	sb.WriteString("]")
	return sb.String()
}
