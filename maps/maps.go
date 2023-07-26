package maps

import "github.com/molikatty/fp"

// Of quickly create a map.
func Of[K comparable, V any](kvs ...fp.Pairs[K, V]) map[K]V {
	var nw = make(map[K]V, len(kvs))
	for i := range kvs {
		nw[kvs[i].Key()] = kvs[i].Value()
	}

	return nw
}

func Make[K comparable, V any](buflen int) map[K]V {
	return make(map[K]V, buflen)
}

// Iter make an iterator for map
func Iter[M ~map[K]V, K comparable, V any](m M) fp.Next[fp.Pairs[K, V]] {
	var yied = make(chan fp.Pairs[K, V])
	go func() {
		for k, v := range m {
			yied <- fp.Pair(k, v)
		}
		close(yied)
	}()

	return func() (fp.Pairs[K, V], bool) {
		for v := range yied {
			return v, true
		}

		return fp.Zero[fp.Pairs[K, V]](), false
	}
}

// IsEmpty check map is empty
func IsEmpty[M ~map[K]V, K comparable, V any](kv M) bool {
	return fp.IsNil(kv) || len(kv) == 0
}

// Clone returns a copy of the map.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[M ~map[K]V, K comparable, V any](kv M) M {
	var nw = make(M, len(kv))
	Copy(kv, nw)

	return nw
}

// Pop randomly pop a key-value pair from the map
func Pop[M ~map[K]V, K comparable, V any](m M) (p fp.Pairs[K, V]) {
	for k, v := range m {
		p = fp.Pair(k, v)
		break
	}

	delete(m, p.Key())
	return
}

// Copy a map into map.
// The elements are copied using assignment, so this is a shallow copy.
func Copy[M ~map[K]V, K comparable, V any](src M, dst M) {
	for k, v := range src {
		dst[k] = v
	}
}

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](kv M) []K {
	var slice = make([]K, 0, len(kv))
	for k := range kv {
		slice = append(slice, k)
	}

	return slice
}

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[M ~map[K]V, K comparable, V any](kv M) []V {
	var slice = make([]V, 0, len(kv))
	for _, v := range kv {
		slice = append(slice, v)
	}

	return slice
}

// Flat multiple maps into one map
func Flat[M ~map[K]V, K comparable, V any](s ...M) M {
	var nw = make(M, fp.Sum(Lens(s...)...))
	for i := range s {
		for k, v := range s[i] {
			nw[k] = v
		}
	}

	return nw
}

// Equal check if two maps are equal.
func Equal[M ~map[K]V, K, V comparable](m1, m2 M) bool {
	return fp.If(len(m1) != len(m2), fp.False, func() bool {
		for k, v1 := range m1 {
			if v2, ok := m2[k]; !ok || v1 != v2 {
				return false
			}
		}
		return true
	})
}

func EqualFunc[M ~map[K]V, K comparable, V any](m1, m2 M, eq func(V1, V2 V) bool) bool {
	return fp.If(len(m1) != len(m2), fp.False, func() bool {
		for k, v1 := range m1 {
			if v2, ok := m2[k]; !ok || !eq(v1, v2) {
				return false
			}
		}
		return true
	})
}

// Lens store the len of multiple slices into a int-slice.
func Lens[M ~map[K]V, K comparable, V any](kvs ...M) []int {
	var n = make([]int, len(kvs))
	for i := range kvs {
		n[i] = len(kvs[i])
	}

	return n
}

// Clear of map
func Clear[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		delete(m, k)
	}
}
