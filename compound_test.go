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
