package compound_test

import (
	"bytes"
	"testing"
	"testing/quick"

	"github.com/tv42/compound"
)

type MyType struct {
	A uint64
	B string
	C uint64
}

func cmp(a, b MyType) int {
	switch {
	case a.A < b.A:
		return -1
	case a.A > b.A:
		return 1
	default:
		switch {
		case a.B < b.B:
			return -1
		case a.B > b.B:
			return 1
		default:
			switch {
			case string(a.C) < string(b.C):
				return -1
			case string(a.C) > string(b.C):
				return 1
			default:
				return 0
			}
		}
	}
}

func TestQuickOrder(t *testing.T) {
	unknown := func(a, b MyType) int {
		as := string(compound.Key(a))
		bs := string(compound.Key(b))
		switch {
		case as < bs:
			return -1
		case as == bs:
			return 0
		case as > bs:
			return 1
		}
		panic("this cannot happen")
	}

	if err := quick.CheckEqual(cmp, unknown, nil); err != nil {
		t.Error(err)
	}
}

func TestQuickRoundtrip(t *testing.T) {
	good := func(a MyType) (MyType, error) {
		return a, nil
	}
	unknown := func(a MyType) (MyType, error) {
		key := compound.Key(a)
		var r MyType
		err := compound.Decode(key, &r)
		return r, err
	}
	if err := quick.CheckEqual(good, unknown, nil); err != nil {
		t.Error(err)
	}
}

func TestQuickErrorShort(t *testing.T) {
	check := func(a MyType) bool {
		key := compound.Key(a)
		for i := 0; i < len(key); i++ {
			var r MyType
			input := key[:i]
			err := compound.Decode(input, &r)
			if err != compound.ErrTooShort {
				t.Errorf("expected too short error for %q, got %v", input, err)
				return false
			}
		}
		return true
	}
	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}
}

func TestQuickPrefixProperty(t *testing.T) {
	check := func(a MyType) bool {
		key := compound.Key(a)
		prefix := compound.PrefixN(a, 2)
		return bytes.HasPrefix(key, prefix)
	}
	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}
}
