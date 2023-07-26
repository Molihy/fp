package slice

// Package slices defines various functions useful with slices of any type.
import (
	"github.com/molikatty/fp"
)

type (
	Cap struct{}
	Len struct{}
)

// Of quickly create a slice.
func Of[T any](t ...T) []T {
	return t
}

func From[S ~[]T, T any](next fp.Next[T]) S {
	return fp.Slice(next)
}

// Iter make an 'iter' for slice
func Iter[S ~[]T, T any](t S) fp.Next[T] {
	var index, stop = -1, len(t)
	return func() (T, bool) {
		index++
		if index == stop {
			return fp.Zero[T](), false
		}

		return t[index], true
	}
}

// Make make an empty slice with bufcap
func Make[S ~[]T, T any](bufcap int) S {
	return make(S, 0, bufcap)
}

// Flat unzip a two-dimensional slice into a one-dimensional slice.
func Flat[S ~[]T, T any](slices ...S) []T {
	var slice = make(S, 0, fp.Sum(Lens(slices...)...))
	for i := range slices {
		slice = append(slice, slices[i]...)
	}

	return slice
}

// Caps store the cap of multiple slices into a int-slice.
func Caps[S ~[]T, T any](slices ...S) []int {
	var caps = make([]int, len(slices))
	for i := range slices {
		caps[i] = cap(slices[i])
	}

	return caps
}

// Lens store the len of multiple slices into a int-slice.
func Lens[S ~[]T, T any](slices ...S) []int {
	var caps = make([]int, len(slices))
	for i := range slices {
		caps[i] = len(slices[i])
	}

	return caps
}

// Trim the unused capacity.
func Trim[S ~[]T, T any](s S) []T {
	return s[:len(s):len(s)]
}

// Buffer return a new slice with the specified buffer capacity (bufcap).
// If bufcap is smaller than the length of the original slice, it will be truncated.
// If bufcap is less than 0, it will cause a panic.
func Buffer[S ~[]T, T any](s S, bufcap int) S {
	if bufcap < 0 {
		panic("cannot be negative")
	}

	return append(make([]T, 0, len(s)), s[:fp.Max(bufcap, len(s))]...)
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[S ~[]T, T any](s S) S {
	return append(make([]T, 0, len(s)), s...)
}

// Clear quickly clear of slice
func Clear[S ~[]T, T any](s S) S {
	return s[:0]
}

// Delete removes the elements s[i:j] from 's', returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
func Delete[S ~[]T, T any](s S, i, j int) S {
	_ = s[i:j]
	return append(s[:i], s[j:]...)
}

// Eq selecting Cap or Len, we can compare the 'cap' or 'len' of two slices.
func Eq[N Cap | Len, S ~[]T, T any](f, s S) bool {
	return fp.If(fp.Is[Len](fp.Zero[N]()),
		func() bool {
			return len(f) == len(s)
		},
		func() bool {
			return cap(f) == cap(s)
		},
	)
}

// Equal compare whether two slices are equal.
func Equal[S ~[]T, T comparable](f, s S) bool {
	if Eq[Len](f, s) {
		return false
	}

	for i := 0; i < len(f); i++ {
		if f[i] != s[i] {
			return false
		}
	}

	return true
}

func EqualFunc[S ~[]T, T any](f, s S, eq func(T, T) bool) bool {
	for i := 0; i < len(s); i++ {
		if !eq(f[i], s[i]) {
			return false
		}
	}

	return true
}

// Index get a value index of slice
func Index[S ~[]T, T comparable](vs S, v T) int {
	for i := range vs {
		if vs[i] == v {
			return i
		}
	}

	return -1
}

func IndexFunc[S ~[]T, T any](vs S, fn func(T) bool) int {
	for i := range vs {
		if fn(vs[i]) {
			return i
		}
	}

	return -1
}

// Contains reports whether v is present in s.
func Contains[S ~[]T, T comparable](vs S, v T) bool {
	return Index(vs, v) >= 0
}

func ContainsFunc[S ~[]T, T any](vs S, fn func(T) bool) bool {
	return IndexFunc(vs, fn) >= 0
}

// Reverse reverses the elements of the slice in place.
func Reverse[S ~[]T, T any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Swap two slice
func Swap[S ~[]T, T any](x, y []T) {
	if len(x) != len(y) {
		panic("x len is not equal to the y len.")
	}

	for i := 0; i < len(x); i++ {
		x[i], y[i] = y[i], x[i]
	}
}
