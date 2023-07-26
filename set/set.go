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
	// Add an element to Set
	Add(K)
	// Adds all elements of a collection to another collection
	Adds(Set[K])
	// Del a element for set
	Del(K)
	// Pop del and return a element arbitrary item from the set
	Pop() K
	// Clear the set
	Clear()
	// IsSafe check the set is concurrency-safe
	IsSafe() bool
	// Has check a element has in set
	Has(K) bool
	// Len return len of set
	Len() int
	// Clone the set
	Clone() Set[K]
	// Union returns a new set with all elements in both sets.
	Union(Set[K]) Set[K]
	// Slice get value for set
	Slice() []K
	// Equal two sets  to each other.
	Equal(Set[K]) bool
	// ForEach the set
	ForEach(func(K) bool)
	// IsEmpty the set
	IsEmpty() bool
	// IsSubset with other set
	IsSubset(Set[K]) bool
	// Intersect return a new set with other set element
	Intersect(Set[K]) Set[K]
	// IsSuperset with other set
	IsSuperset(Set[K]) bool
	// Difference return a new set with other set element
	Difference(Set[K]) Set[K]
	// IsProperSubset with other set
	IsProperSubset(Set[K]) bool
	// IsProperSuperset with other set
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

func From[U Safe | Unsafe, K comparable](next fp.Next[K]) Set[K] {
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
