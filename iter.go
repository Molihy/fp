package fp

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// Next iterator type with some processing functions implemented for it.
type Next[E any] func() (E, bool)

func (next Next[E]) String() string {
	var s = strings.Split(runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name(), "/")
	return fmt.Sprintf("<%v %p>", s[len(s)-1], next)
}

// Filter the iterator, keeping only the values that evaluate to 'true'.
// this function is lazy.
func Filter[E any](next Next[E], fn func(E) bool) Next[E] {
	return func() (E, bool) {
		for {
			t, ok := next()
			if Not(ok) {
				return Zero[E](), false
			}

			if fn(t) {
				return t, true
			}
		}
	}
}

// Reduce from left to right, 1, 2, 3, 4 will be computed as (((1+2)+3)+4),
// resulting in a single value. the function is not lazy
func Reduce[E any](next Next[E], fn func(E, E) E) E {
	var poly E
	for {
		t, ok := next()
		if Not(ok) {
			return poly
		}

		poly = fn(poly, t)
	}
}

// Fold from left to right, 1, 2, 3, 4 will be computed as 1, (1+2), ((1+2)+3), (((1+2)+3)+4)
func Fold[E any](next Next[E], fn func(E, E) E) Next[E] {
	var poly E
	return func() (E, bool) {
		e, ok := next()
		if Not(ok) {
			return Zero[E](), false
		}

		poly = fn(poly, e)
		return poly, true
	}
}

// Map the values inside the iterator to new values
// this function is lazy.
func Map[ER, E any](next Next[E], fn func(E) ER) Next[ER] {
	return func() (ER, bool) {
		t, ok := next()
		if Not(ok) {
			return Zero[ER](), false
		}

		return fn(t), true
	}
}

// Zip the values at the beginning of multiple iterators.
// this function is lazy.
func Zip[E any](nexts ...Next[E]) Next[[]E] {
	return func() ([]E, bool) {
		var slice = make([]E, 0)
		for i := range nexts {
			t, ok := nexts[i]()
			if Not(ok) {
				return Zero[[]E](), false
			}
			slice = append(slice, t)
		}

		return slice, true
	}
}

// Lock the iterator to make it concurrent-safe.
// this function is lazy.
func Lock[E any](next Next[E]) Next[E] {
	var lock sync.Mutex
	return func() (E, bool) {
		lock.Lock()
		defer lock.Unlock()

		return next()
	}
}

// Take iterate the iterator only for a specific number of times
// this function is lazy.
func Take[E any](stop int, next Next[E]) Next[E] {
	var n int
	return func() (E, bool) {
		n++
		if n > stop {
			return Zero[E](), false
		}

		return next()
	}
}

// ForEach loop through an iterator to retrieve values, and you can stop the loop prematurely by
// returning false from the 'fn'.
func ForEach[E any](next Next[E], fn func(E) bool) {
	for n, ok := next(); ok; n, ok = next() {
		if Not(fn(n)) {
			return
		}
	}
}

// Stop add a stopping condition to the iterator, often used to stop infinite iterators.
// this function is lazy.
func Stop[E any](next Next[E], fn func(E) bool) Next[E] {
	return func() (E, bool) {
		t, _ := next()
		if fn(t) {
			return Zero[E](), false
		}

		return t, true
	}
}

// From make the iterator compatible with the 'range'.
// When the iterator is empty, the channel will be closed.
//
//	for v := range From(Iter){
//	  // do something
//	}
func From[E any](next Next[E]) <-chan E {
	var y = make(chan E)
	go func() {
		for n, ok := next(); ok; n, ok = next() {
			y <- n
		}

		close(y)
	}()

	return y
}

// Yield iteration in Go
func Yield[E any](next Next[E]) (yield func() E) {
	var y = From(next)
	return func() E {
		return <-y
	}
}

// Range generate a range of integers, similar to Python built-in 'range' function.
func Range[N Integer](r ...N) Next[N] {
	var iter = func(start, stop, step N) Next[N] {
		return func() (N, bool) {
			start += step
			if start >= stop {
				return Zero[N](), false
			}

			return start, true
		}
	}

	switch len(r) {
	case 0:
		panic(ErrLeastOne)
	case 1:
		return iter(Zero[N]()-1, r[0], 1)
	case 2:
		return iter(r[0]-1, r[1], 1)
	default:
		return iter(r[0]-r[2], r[1], r[2])
	}
}

// Slice generate a slice from an iterator.
func Slice[E any](next Next[E]) []E {
	var slice = make([]E, 0)
	ForEach(next, func(e E) bool {
		slice = append(slice, e)
		return true
	})

	return slice
}

// Chan generate a channel with a custom bufcap from an iterator.
func Chan[E any](next Next[E], bufcap int) chan E {
	var channel = make(chan E, bufcap)
	go func() {
		ForEach(next, func(t E) bool {
			channel <- t
			return true
		})
		close(channel)
	}()

	return channel
}

// KV generate a map from iterator
func KV[K comparable, V any](next Next[Pairs[K, V]]) map[K]V {
	var m = make(map[K]V)
	ForEach(next, func(p Pairs[K, V]) bool {
		m[p.Key()] = p.Value()
		return true
	})

	return m
}

// Iota infinite iterator for numbers, lazily generates an infinite sequence of numbers
//
//	Iota() -> 0 ... N
//	Iota(10) -> 10 ... N
//	Iota(5, 3) -> 5 -> 8 -> 11 -> ... N
func Iota[N Number](n ...N) Next[N] {
	var iter = func(start, step N) Next[N] {
		return func() (N, bool) {
			start += step
			return start, true
		}
	}

	switch len(n) {
	case 0:
		return iter(Zero[N]()-1, 1)
	case 1:
		return iter(n[0]-1, 1)
	default:
		return iter(n[0]-n[1], n[1])
	}
}

// All check if there is any 'false' in the bool iterator, similar to Python
// built-in function 'all'.
func All(next Next[bool]) (ok bool) {
	ForEach(next, func(b bool) bool {
		ok = b
		return b
	})

	return
}

// Merge multiple iterators, will iterate in order from left to right
func Merge[E any](nexts ...Next[E]) Next[E] {
	var index, length = 0, len(nexts)
	var next Next[E]
	next = func() (E, bool) {
		if index == length {
			return Zero[E](), false
		}

		e, ok := nexts[index]()
		if Not(ok) {
			index++
			return next()
		}

		return e, true
	}

	return next
}
