package compound_test

import (
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
