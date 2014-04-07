package compound_test

import (
	"testing"

	"github.com/tv42/compound"
)

func TestKeyString(t *testing.T) {
	type T struct {
		S string
	}
	val := T{S: "foo"}
	key := compound.Key(val)
	if g, e := string(key), "\x02foo\x00"; g != e {
		t.Errorf("bad Key: %#v -> %q != %q", val, g, e)
	}
}

func TestKeyStringEmpty(t *testing.T) {
	type T struct {
		S string
	}
	val := T{S: ""}
	key := compound.Key(val)
	if g, e := string(key), "\x02\x00"; g != e {
		t.Errorf("bad Key: %#v -> %q != %q", val, g, e)
	}
}

func TestKeyStringNul(t *testing.T) {
	type T struct {
		S string
	}
	val := T{S: "foo\x00bar"}
	key := compound.Key(val)
	if g, e := string(key), "\x02foo\x00\xFFbar\x00"; g != e {
		t.Errorf("bad Key: %#v -> %q != %q", val, g, e)
	}
}

func TestKeyBytes(t *testing.T) {
	type T struct {
		S []byte
	}
	val := T{S: []byte("foo")}
	key := compound.Key(val)
	if g, e := string(key), "\x01foo\x00"; g != e {
		t.Errorf("bad Key: %#v -> %q != %q", val, g, e)
	}
}

func TestKeyUint64(t *testing.T) {
	type T struct {
		N uint64
	}
	val := T{N: 532}
	key := compound.Key(val)
	if g, e := string(key), "\x03\x00\x00\x00\x00\x00\x00\x02\x14"; g != e {
		t.Errorf("bad Key: %#v -> %q != %q", val, g, e)
	}
}

func TestDecodeString(t *testing.T) {
	type T struct {
		S string
	}
	key := []byte("\x02foo\x00")
	var val T
	if err := compound.Decode(key, &val); err != nil {
		t.Fatal(err)
	}
	if g, e := val.S, "foo"; g != e {
		t.Errorf("bad Decode: %#v -> %q != %q", key, g, e)
	}
}

func TestPrefixN(t *testing.T) {
	type K struct {
		A uint64
		B string
		C uint64
	}
	val := K{A: 42, B: "noise"}
	fields := 1
	key := compound.PrefixN(val, fields)
	if g, e := string(key), "\x03\x00\x00\x00\x00\x00\x00\x00*"; g != e {
		t.Errorf("bad PrefixN: %#v %d -> %q != %q", val, fields, g, e)
	}
}
