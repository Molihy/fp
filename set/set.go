package set

import (
	"github.com/molikatty/fp"
	"github.com/molikatty/fp/channel"
)

type (
	Safe   struct{}
	Unsafe struct{}
)

type Set[K comparable] interface {
	Add(K)
	Adds(Set[K])
	Del(K)
	Pop() K
	Clear()
	IsSafe() bool
	Has(K) bool
	Len() int
	Clone() Set[K]
	Union(Set[K]) Set[K]
	Slice() []K
	Equal(Set[K]) bool
	ForEach(func(K) bool)
	IsEmpty() bool
	IsSubset(Set[K]) bool
	Intersect(Set[K]) Set[K]
	IsSuperset(Set[K]) bool
	Difference(Set[K]) Set[K]
	IsProperSubset(Set[K]) bool
	IsProperSuperset(Set[K]) bool
}

func Of[U Safe | Unsafe, K comparable](t ...K) Set[K] {
	var safe = func() Set[K] {
		return safeOf(t...)
	}

	var unsafe = func() Set[K] {
		return unsafeOf(t...)
	}

	return fp.If(fp.Is[Safe](fp.Zero[U]()), safe, unsafe)
}

func OfFrom[U Safe | Unsafe, K comparable](next fp.Next[K]) Set[K] {
	var s = Of[U, K]()
	fp.ForEach(next, func(k K) bool {
		s.Add(k)
		return true
	})

	return s
}

func Iter[K comparable](set Set[K]) fp.Next[K] {
	var y = channel.Lazy(make(chan K), func(y chan K) {
		set.ForEach(func(k K) bool {
			y <- k
			return true
		})
	})

	return func() (K, bool) {
		for v := range y {
			return v, true
		}

		return fp.Zero[K](), false
	}
}
