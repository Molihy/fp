package set

import (
	"fmt"

	"github.com/molikatty/fp"
	"github.com/molikatty/fp/str"
)

type _unsafe[K comparable] map[K]fp.None

var _ Set[struct{}] = fp.Zero[_unsafe[struct{}]]()

func unsafeOf[K comparable](t ...K) _unsafe[K] {
	var set = make(_unsafe[K], len(t))
	for i := range t {
		set[t[i]] = fp.Zero[fp.None]()
	}

	return set
}

func (s _unsafe[K]) Add(t K) {
	s[t] = fp.Zero[fp.None]()
}

func (s _unsafe[K]) Adds(other Set[K]) {
	other.Loop(func(k K) { s.Add(k) })
}

func (s _unsafe[K]) Del(t K) {
	delete(s, t)
}

func (s _unsafe[K]) Len() int {
	return len(s)
}

func (s _unsafe[K]) Pop() (t K) {
	s.ForEach(func(k K) bool {
		t = k
		s.Del(k)
		return false
	})

	return
}

func (s _unsafe[K]) IsSafe() bool {
	return false
}

func (s _unsafe[K]) Clone() Set[K] {
	var newSet = unsafeOf[K]()
	newSet.Adds(s)
	return newSet
}

func (s _unsafe[K]) IsEmpty() bool {
	return fp.IsNil(s) || s.Len() == 0
}

func (s _unsafe[K]) Has(t K) bool {
	return fp.InMap(s, t)
}

func (s _unsafe[K]) Clear() {
	s.Loop(func(k K) { s.Del(k) })
}

func (s _unsafe[K]) ForEach(fn func(K) bool) {
	for k := range s {
		if !fn(k) {
			return
		}
	}
}

func (s _unsafe[K]) Loop(fn func(K)) {
	for k := range s {
		fn(k)
	}
}

func (s _unsafe[K]) Union(other Set[K]) Set[K] {
	var newSet = make(_unsafe[K], fp.Max(s.Len(), other.Len()))
	newSet.Adds(s)
	newSet.Adds(other)

	return newSet
}

func (s _unsafe[K]) String() string {
	var strs = fp.Slice(fp.Map[string](Iter[K](s), func(k K) string {
		return fmt.Sprint(k)
	}))

	return "{" + str.Join(", ", strs...) + "}"
}

func (s _unsafe[K]) Slice() []K {
	var slice = make([]K, 0, s.Len())
	s.Loop(func(k K) { slice = append(slice, k) })

	return slice
}

func (s _unsafe[K]) Equal(other Set[K]) bool {
	var check = func() (b bool) {
		s.ForEach(func(k K) bool {
			b = other.Has(k)
			return b
		})

		return
	}

	return fp.If(s.Len() != other.Len(), fp.False, check)
}

func (s _unsafe[K]) Difference(other Set[K]) Set[K] {
	var newSet = unsafeOf[K]()
	s.Loop(func(k K) {
		if !other.Has(k) {
			newSet.Add(k)
		}
	})

	return newSet
}

func (s _unsafe[K]) Intersect(other Set[K]) Set[K] {
	var newSet = unsafeOf[K]()
	var t = fp.If(s.Len() < other.Len(),
		func() fp.Pairs[Set[K], Set[K]] {
			return fp.Pair(fp.AnyTo[Set[K]](s), other)
		},
		func() fp.Pairs[Set[K], Set[K]] {
			return fp.Pair(other, fp.AnyTo[Set[K]](s))
		},
	)

	t.Key().Loop(func(k K) {
		if t.Value().Has(k) {
			newSet.Add(k)
		}
	})

	return newSet
}

func (s _unsafe[K]) IsSubset(other Set[K]) (ok bool) {
	if s.Len() > other.Len() {
		return false
	}

	s.ForEach(func(k K) bool {
		ok = fp.If(fp.Not(other.Has(k)), fp.False, fp.True)
		return ok
	})

	return
}

func (s _unsafe[K]) IsSuperset(other Set[K]) bool {
	return other.IsSubset(s)
}

func (s _unsafe[K]) IsProperSubset(other Set[K]) bool {
	return s.IsSubset(other) && s.Len() != other.Len()
}

func (s _unsafe[K]) IsProperSuperset(other Set[K]) bool {
	return s.IsSuperset(other) && s.Len() != other.Len()
}
