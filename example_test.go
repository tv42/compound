package compound_test

import (
	"bytes"
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

func ExamplePrefixN() {
	type MyKey struct {
		User    string
		Album   string
		PhotoID uint64
	}
	val := MyKey{User: "jdoe", Album: "Travel", PhotoID: 42}
	fmt.Printf("%q\n", compound.PrefixN(val, 1))
	fmt.Printf("%q\n", compound.PrefixN(val, 2))
	fmt.Printf("%q\n", compound.PrefixN(val, 3))
	// Output:
	// "\x02jdoe\x00"
	// "\x02jdoe\x00\x02Travel\x00"
	// "\x02jdoe\x00\x02Travel\x00\x03\x00\x00\x00\x00\x00\x00\x00*"
}

func ExamplePrefixNPartial() {
	type MyKey struct {
		User    string
		Album   string
		PhotoID uint64
	}

	item := MyKey{User: "jdoe", Album: "Travel in Europe", PhotoID: 42}
	key := compound.Key(item)
	fmt.Printf("%q\n", key)

	query := MyKey{User: "jdoe", Album: "Travel in "}
	prefix := compound.PrefixNPartial(query, 2)
	fmt.Printf("%q\n", prefix)

	fmt.Printf("matches: %v\n", bytes.HasPrefix(key, prefix))
	// Output:
	// "\x02jdoe\x00\x02Travel in Europe\x00\x03\x00\x00\x00\x00\x00\x00\x00*"
	// "\x02jdoe\x00\x02Travel in "
	// matches: true
}
