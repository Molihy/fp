package fp

import (
	"context"
	"unsafe"
)

type (
	FAny   = func(any) any
	FAnys  = func(...any) any
	None   = struct{}
	Run    = func(func())
	RunCtx = func(func(context.Context))
	Wait   = func()

	// Generic collection of signed numbers
	Signed interface {
		~int8 | ~int16 | ~int32 | ~int64 | ~int
	}

	// Generic collection of unsigned numbers
	Unsigned interface {
		~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr
	}

	// Generic collection of integers
	Integer interface {
		Signed | Unsigned
	}

	// Generic collection of float numbers
	Float interface {
		~float32 | ~float64
	}

	// Generic collection of all numeric types
	Number interface {
		Integer | Float | ~complex64 | ~complex128
	}

	// Generic collection of data types with size
	Size interface {
		Float | Signed | Unsigned | string
	}
)

// To can unsafe convert one type to another type. If the type is interface{},
// it will be obtained using a type assertion. If the assertion fails, it will
// be cast to another type.
//
// Warning: If you want to convert string to the []rune or []byte, please use str.To().
func To[R, T any](a T) R {
	return *(*R)(unsafe.Pointer(&a))
}

// Sizeof get the size of a type
func Sizeof[T any](a T) int {
	return To[int](unsafe.Sizeof(a))
}

// Any convert any type to an empty interface
func Any[T any](a T) any {
	return a
}

// AnyTo convert a value of type 'any' to the specified type safely
func AnyTo[T any](a any) T {
	return a.(T)
}

// ToFAny is func(A) B conver to FAny
func ToFAny[A, B any](f func(A) B) FAny {
	return func(a any) any {
		return f(AnyTo[A](a))
	}
}

// FAnyTo is FAny conver to func(A) B
func FAnyTo[A, B any](fn func(any) any) func(A) B {
	return func(a A) B {
		return AnyTo[B](fn(a))
	}
}
