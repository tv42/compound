package compound_test

import (
	"testing"

	"github.com/tv42/compound"
)

func BenchmarkKey(b *testing.B) {
	type MyKey struct {
		User    string
		Album   string
		PhotoID uint64
	}
	val := MyKey{User: "jdoe", Album: "Travel", PhotoID: 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := compound.Key(val)
		_ = key
	}
}

func BenchmarkDecode(b *testing.B) {
	type MyKey struct {
		User    string
		Album   string
		PhotoID uint64
	}
	key := compound.Key(MyKey{User: "jdoe", Album: "Travel", PhotoID: 42})
	var val MyKey
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := compound.Decode(key, &val)
		if err != nil {
			b.Error(err)
		}
	}
}
