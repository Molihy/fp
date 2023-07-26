package set

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/molikatty/fp"
	"github.com/molikatty/fp/str"
)

type _safe[K comparable] struct {
	set sync.Map
	len atomic.Int64
}

var _ Set[struct{}] = fp.Of[_safe[struct{}]]()

func safeOf[K comparable](t ...K) *_safe[K] {
	var s = fp.Of[_safe[K]]()
	for i := range t {
		s.set.Store(t[i], fp.Zero[fp.None]())
	}
	s.len.Store(fp.To[int64](len(t)))

	return s
}

func (s *_safe[K]) Add(t K) {
	s.set.Store(t, fp.Zero[fp.None]())
	s.len.Add(1)
}

func (s *_safe[K]) Adds(other Set[K]) {
	other.ForEach(func(k K) bool {
		s.Add(k)
		return true
	})

}

func (s *_safe[K]) Del(t K) {
	_, ok := s.set.LoadAndDelete(t)
	if ok {
		s.len.Add(-1)
	}
}

func (s *_safe[K]) Clone() Set[K] {
	var newSet = safeOf[K]()
	newSet.Adds(s)
	return newSet
}

func (s *_safe[K]) IsSafe() bool {
	return true
}

func (s *_safe[K]) IsEmpty() bool {
	return fp.IsNil(s) || s.Len() == 0
}

func (s *_safe[K]) Len() int {
	return fp.To[int](s.len.Load())
}

func (s *_safe[K]) Pop() (t K) {
	s.ForEach(func(k K) bool {
		pop, _ := s.set.LoadAndDelete(k)
		t = fp.AnyTo[K](pop)
		return false
	})

	return
}

func (s *_safe[K]) Has(t K) bool {
	_, ok := s.set.Load(t)
	return ok
}

func (s *_safe[K]) Clear() {
	s.ForEach(func(k K) bool {
		s.Del(k)
		return true
	})
}

func (s *_safe[K]) Union(other Set[K]) Set[K] {
	var newSet = safeOf[K]()

	s.ForEach(func(k K) bool {
		newSet.Add(k)
		return true
	})

	other.ForEach(func(k K) bool {
		newSet.Add(k)
		return true
	})

	return newSet
}

func (s *_safe[K]) ForEach(fn func(K) bool) {
	s.set.Range(func(k, _ any) bool {
		return fn(fp.AnyTo[K](k))
	})
}

func (s *_safe[K]) String() string {
	var strs = fp.Slice(fp.Map[string](Iter[K](s), func(k K) string {
		return fmt.Sprint(k)
	}))

	return "{" + str.Join(", ", strs...) + "}"
}

func (s *_safe[K]) Slice() []K {
	var slice = make([]K, 0, s.Len())
	s.ForEach(func(k K) bool {
		slice = append(slice, k)
		return true
	})

	return slice
}

func (s *_safe[K]) Equal(other Set[K]) bool {
	var check = func() (b bool) {
		s.ForEach(func(k K) bool {
			b = other.Has(k)
			return b
		})

		return
	}

	return fp.If(s.Len() != other.Len(), fp.False, check)
}

func (s *_safe[K]) Difference(other Set[K]) Set[K] {
	var newSet = safeOf[K]()
	s.ForEach(func(k K) bool {
		if !other.Has(k) {
			newSet.Add(k)
		}

		return true
	})

	return newSet
}

func (s *_safe[K]) Intersect(other Set[K]) Set[K] {
	var newSet = safeOf[K]()
	var t = fp.If(s.Len() < other.Len(),
		func() fp.Pairs[Set[K], Set[K]] {
			return fp.Pair(fp.To[Set[K]](s), other)
		},
		func() fp.Pairs[Set[K], Set[K]] {
			return fp.Pair(other, fp.To[Set[K]](s))
		},
	)

	t.Key().ForEach(func(k K) bool {
		if t.Value().Has(k) {
			newSet.Add(k)
		}

		return true
	})

	return newSet
}

func (s *_safe[K]) IsSubset(other Set[K]) (ok bool) {
	if s.Len() > other.Len() {
		return false
	}

	s.ForEach(func(k K) bool {
		ok = fp.If(fp.Not(other.Has(k)), fp.False, fp.True)
		return ok
	})

	return
}

func (s *_safe[K]) IsSuperset(other Set[K]) bool {
	return other.IsSubset(s)
}

func (s *_safe[K]) IsProperSubset(other Set[K]) bool {
	return s.IsSubset(other) && s.Len() != other.Len()
}

func (s *_safe[K]) IsProperSuperset(other Set[K]) bool {
	return s.IsSuperset(other) && s.Len() != other.Len()
}
