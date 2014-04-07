# compound keys as `[]byte`

`compound` is a Go library that helps you make order-preserving keys
from multiple items of data, with support for prefix iteration. It
does this even with variable length items of arbitrary data (of the
supported data types).

It works like this: first, define your struct

``` go
type MyKey struct {
	User string
	Album string
	PhotoID uint64
}
```

Generate a key for a record:

``` go
var key []byte
key = compound.Key(MyKey{User: "jdoe", Album: "Travel", PhotoID: 42})
```

You can decode a key:

``` go
var data MyKey
err := compound.Decode(key, &data)
```

You can ask for a prefix that matches all keys with specific values:

``` go
// iterate an album
prefix := compound.Prefix(MyKey{User: "jdoe", Album: "Travel"})
```

The above works by skipping zero valued items at the end. Sometimes
the zero value is necessary, so you can also be more explicit:

``` go
prefix := compound.PrefixN(MyKey{User: "jdoe", Album: "Travel"}, 2)
```

These match all keys with the keys at those exact values. For example,
the above will not match an album named "Travel in Europe".

To ask for matching partial values of the last key, use:

``` go
prefix := compound.PrefixPartial(MyKey{User: "jdoe", Album: "Travel in "})
prefix := compound.PrefixNPartial(MyKey{User: "jdoe", Album: "Travel in "}, 2)
```

This does not terminate the last key, letting the above match "Travel
in Europe".

Prefixes can also be decoded:

``` go
// for keys from Prefix or PrefixN
var data MyKey
err := compound.DecodePrefix(key, &data)
```

``` go
// for keys from PrefixPartial or PrefixNPartial
var data MyKey
err := compound.DecodePrefixPartial(key, &data)
```
