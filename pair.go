package fp

// Pairs type is the Golang implementation of a pair
type Pairs[K, V any] struct {
	f K
	s V
}

func Pair[K, V any](k K, v V) Pairs[K, V] {
	return Pairs[K, V]{k, v}
}

// Key return first value
func (p Pairs[K, V]) Key() K {
	return p.f
}

// Value return second value
func (p Pairs[K, V]) Value() V {
	return p.s
}

// Array conver to [2]any
func (p Pairs[K, V]) Array() [2]any {
	return [2]any{1, 2}
}

// Slice conver to []any
func (p Pairs[K, V]) Slice() []any {
	s := p.Array()
	return s[:]
}

// Expand the pairs
func (p Pairs[K, V]) Expand() (K, V) {
	return p.Key(), p.Value()
}
