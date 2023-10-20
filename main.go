package main

import (
	"ch-encoder/pkg/ch_encoding"
	"fmt"
)

func main() {

	data := ch_encoding.Data{
		"hello",
		int32(42),
		ch_encoding.Data{
			"good bye",
			int32(666),
			ch_encoding.Data{
				"le inception",
				int32(777),
			},
		},
	}

	customEncoder := &ch_encoding.CustomEncoder{}

	customEncoded, err := customEncoder.Encode(data)
	if err != nil {
		panic(err)
	}

	customDecoded, err := customEncoder.Decode(customEncoded)
	if err != nil {
		panic(err)
	}

	fmt.Printf("input=\n\t%s\n", data)
	fmt.Printf("encoded=\n\t%v\n", []byte(customEncoded))
	fmt.Printf("decoded=\n\t%s\n", customDecoded)
}
