// Package compound creates order-preserving binary keys from
// structured data.
//
// The generated keys support prefix iteration, and may contain
// variable length items of arbitrary data (of the supported data
// types).
package compound

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

const (
	codeReservedZero byte = 0x00
	codeBytes        byte = 0x01
	codeString       byte = 0x02
	codeUint64       byte = 0x03
	// must be greater than all the other possible values
	codeReservedQuote byte = 0xFF
)

func quote(key []byte, data []byte) []byte {
	data = bytes.Replace(data, []byte{0x00}, []byte{0x00, codeReservedQuote}, -1)
	key = append(key, data...)
	key = append(key, 0x00)
	return key
}

// Key returns bytes that encode the key in a non-ambiguous,
// order-preserving way. The key can later be decoded to recover the
// individual fields.
//
// Key panics if value is not a struct or a pointer to a struct,
// or if a field is of an unsupported type.
func Key(value interface{}) []byte {
	var key []byte
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("compound.Key only works with structs: %T: %#v", value, value))
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		i := f.Interface()
		switch x := i.(type) {
		case []byte:
			key = append(key, codeBytes)
			key = quote(key, x)
		case string:
			key = append(key, codeString)
			key = quote(key, []byte(x))
		case uint64:
			key = append(key, codeUint64)
			var b [8]byte
			binary.BigEndian.PutUint64(b[:], x)
			key = append(key, b[:]...)
		default:
			panic(fmt.Errorf("compound.Key unsupported field type: %T: %#v", i, i))
		}

	}
	return key
}

// func Prefix(value interface{}) []byte {
// }

// func PrefixN(value interface{}, n int) []byte {
// }

// func PrefixPartial(value interface{}) []byte {
// }

// func PrefixNPartial(value interface{}, n int) []byte {
// }

// func Decode(key []byte, ptr interface{}) error {
// }

// func DecodePrefix(key []byte, ptr interface{}) error {
// }

// func DecodePrefixPartial(key []byte, ptr interface{}) error {
// }
