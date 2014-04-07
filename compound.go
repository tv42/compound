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
	"errors"
	"fmt"
	"reflect"
)

// ErrTooShort is returned when decoding a key that does not contain
// the end marker.
var ErrTooShort = errors.New("key is too short")

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

func unquote(key []byte) (n int, val []byte, err error) {
	for {
		i := bytes.IndexByte(key, 0x00)
		if i == -1 {
			return n, nil, ErrTooShort
		}
		n += i
		val = append(val, key[:i+1]...)
		switch {
		// cannot underflow because i would have been -1 above
		case i == len(key)-1:
			// looks like the last key, that's ok

		// it is now safe to peek one past i
		case key[i+1] == codeReservedQuote:
			// it's a quoted nul; we already put it in val above
			key = key[i:]
			continue

		default:
			// it's not a quoted nul, so bail out
		}

		// take back the extra nul we put in val above
		val = val[:len(val)-1]
		// mark the nul consumed
		n += 1
		return
	}
}

func checkType(value interface{}) reflect.Value {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("compound.Key only works with structs: %T: %#v", value, value))
	}
	return v
}

func makeKey(v reflect.Value, n int) []byte {
	var key []byte
	for i := 0; i < n; i++ {
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

// Key returns bytes that encode the key in a non-ambiguous,
// order-preserving way. The key can later be decoded to recover the
// individual fields.
//
// Key panics if value is not a struct or a pointer to a struct,
// or if a field is of an unsupported type.
func Key(value interface{}) []byte {
	v := checkType(value)
	return makeKey(v, v.NumField())
}

// func Prefix(value interface{}) []byte {
// }

// PrefixN encodes the first n fields. See Key for more.
func PrefixN(value interface{}, n int) []byte {
	v := checkType(value)
	return makeKey(v, n)
}

// func PrefixPartial(value interface{}) []byte {
// }

// PrefixNPartial encodes the first n fields. The last key will not be
// marked as complete, allowing to perform prefix matching within the
// value. See Key for more.
func PrefixNPartial(value interface{}, n int) []byte {
	v := checkType(value)
	k := makeKey(v, n)
	k = k[:len(k)-1]
	return k
}

// Decode fills the struct at ptr with the values of the fields in
// key.
//
// ptr must be a pointer to a value of the same type as was passed to
// Key earlier.
//
// Decode panics if ptr is not a pointer to a struct,
// or if a field is of an unsupported type.
func Decode(key []byte, ptr interface{}) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Errorf("compound.Decode only works with pointers to structs: %T: %#v", ptr, ptr))
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("compound.Decode only works with pointers to structs: %T: %#v", ptr, ptr))
	}

	for i := 0; i < v.NumField(); i++ {
		if len(key) == 0 {
			return ErrTooShort
		}
		code := key[0]
		key = key[1:]
		f := v.Field(i)
		i := f.Interface()
		switch i.(type) {
		case []byte:
			if code != codeBytes {
				return fmt.Errorf("corrupt key: expecting bytes, got code %x", code)
			}
			n, b, err := unquote(key)
			if err != nil {
				return err
			}
			key = key[n:]
			f.SetBytes(b)
		case string:
			if code != codeString {
				return fmt.Errorf("corrupt key: expecting string, got code %x", code)
			}
			n, b, err := unquote(key)
			if err != nil {
				return err
			}
			key = key[n:]
			f.SetString(string(b))
		case uint64:
			if code != codeUint64 {
				return fmt.Errorf("corrupt key: expecting uint64, got code %x", code)
			}
			if len(key) < 8 {
				return ErrTooShort
			}
			f.SetUint(binary.BigEndian.Uint64(key[:8]))
			key = key[8:]
		default:
			panic(fmt.Errorf("compound.Key unsupported field type: %T: %#v", i, i))
		}

	}
	return nil
}

// func DecodePrefix(key []byte, ptr interface{}) error {
// }

// func DecodePrefixPartial(key []byte, ptr interface{}) error {
// }
