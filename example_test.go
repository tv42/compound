package compound_test

import (
	"fmt"

	"github.com/tv42/compound"
)

func ExampleKey() {
	type MyKey struct {
		User    string
		Album   string
		PhotoID uint64
	}
	var key []byte
	key = compound.Key(MyKey{User: "jdoe", Album: "Travel", PhotoID: 42})
	fmt.Printf("%q\n", key)
	// Output:
	// "\x02jdoe\x00\x02Travel\x00\x03\x00\x00\x00\x00\x00\x00\x00*"
}

func ExampleDecode() {
	type MyKey struct {
		User    string
		Album   string
		PhotoID uint64
	}
	var key = []byte("\x02jdoe\x00\x02Travel\x00\x03\x00\x00\x00\x00\x00\x00\x00*")
	var val MyKey
	if err := compound.Decode(key, &val); err != nil {
		fmt.Println("handle errors:", err)
		return
	}
	fmt.Printf("%#v\n", val)
	// Output:
	// compound_test.MyKey{User:"jdoe", Album:"Travel", PhotoID:0x2a}
}
