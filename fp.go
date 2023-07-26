package fp

import (
	"context"
	"reflect"
	"sync"
)

// Curry is currying of any func
func Curry(fn any) FAnys {
	var (
		fv    = reflect.ValueOf(fn)
		n     = fv.Type().NumIn()
		curry FAnys
	)

	curry = func(args ...any) any {
		var result = func() any {
			return apply(fv, args)
		}

		var again = func() any {
			return func(margs ...any) any {
				return curry(append(args, margs...)...)
			}
		}

		return If(len(args) >= n, result, again)
	}

	return curry
}

// Compose is combines functions from left to right
func Compose[T any](fs ...func(T) T) func(T) T {
	var again = func(a T) T {
		return fs[0](Compose(fs[1:]...)(a))
	}

	return If(IsZero(len(fs)), Lazy(Id[T]), Lazy(again))
}

// Pipe is combines functions from right to left
func Pipe[T any](fs ...func(T) T) func(T) T {
	var again = func(a T) T {
		return fs[len(fs)-1](Compose(fs[:len(fs)-1]...)(a))
	}

	return If(IsZero(len(fs)), Lazy(Id[T]), Lazy(again))
}

// Memoize is caching the return value of a function
func Memoize[T comparable, R any](fn func(T) R) func(T) R {
	var memoize = make(map[T]R)
	return func(t T) R {
		has := func() R {
			return memoize[t]
		}

		nohash := func() R {
			memoize[t] = fn(t)
			return memoize[t]
		}

		return If(Has(memoize, t), has, nohash)
	}
}

// MemoizeFunc is Caching the return value of a function,
// fc transforms the formal parameters by converting an
// incomparable formal parameter into a comparable one, and then caches it.
func MemoizeFunc[T comparable, A, R any](fn func(A) R, fc func(A) T) func(A) R {
	var memoize = make(map[T]R)
	return func(a A) R {
		// A type conver to T type
		key := fc(a)
		has := func() R {
			return memoize[key]
		}

		nohash := func() R {
			memoize[key] = fn(a)
			return memoize[key]
		}

		return If(Has(memoize, key), has, nohash)
	}
}

// Apply is like nodejs function apple
func Apply(fn any, args []any) any {
	var fv = reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		panic("can't Apple a non-function type")
	}

	return apply(fv, args)
}

func apply(fv reflect.Value, args []any) any {
	var argvs = make([]reflect.Value, 0, len(args))
	for i := range args {
		argvs = append(argvs, reflect.ValueOf(args[i]))
	}

	var answer = fv.Call(argvs)
	switch len(answer) {
	case 0:
		return nil
	case 1:
		return answer[0].Interface()
	default:
		anys := make([]any, 0, len(answer))
		for i := range answer {
			anys = append(anys, answer[i].Interface())
		}

		return anys
	}
}

// Ptr used to indirectly obtain a pointer to data that
// cannot be directly obtained as a pointer. example obtaining pointers to
// literal data and function return values
//
//	p1 := Ptr(100)
//	p2 := Ptr(Foo())
//
// Warning: This function cannot obtain a pointer to a variable. If you use this
// function to obtain a pointer to a variable, you will only get a pointer to a copy
// of the variable's value.
func Ptr[T any](t T) *T {
	return &t
}

// Of get a type pointer Instead of using the new() function to generate pointers,
// it aligns with the functional thinking.
//
// Warning: will cause the pointer to escape to the heap, resulting in some
// performance overhead.
func Of[T any]() *T {
	return new(T)
}

// Zero get zero value of a type
func Zero[T any]() (zero T) {
	return
}

// Elem get a value of pointer
func Elem[T any](t *T) T {
	return *t
}

// Min smallest of values.
func Min[T Size](n ...T) T {
	if len(n) < 1 {
		panic(ErrLeastOne)
	}

	var min T
	for i := range n {
		min = If(n[i] < min, Lazy(n[i]), Lazy(min))
	}

	return min
}

// MinFrom get smallest of iterated.
func MinFrom[T Size](n Next[T]) T {
	var min T
	for m := range From(n) {
		min = If(m < min, Lazy(m), Lazy(min))
	}

	return min
}

// Max get maximum of values
func Max[T Size](n ...T) T {
	if len(n) < 1 {
		panic(ErrLeastOne)
	}

	var max T
	for i := range n {
		max = If(n[i] > max, Lazy(n[i]), Lazy(max))
	}

	return max
}

// MaxFrom get maximum of iterated
func MaxFrom[T Size](n Next[T]) T {
	var max T
	for m := range From(n) {
		max = If(m > max, Lazy(m), Lazy(max))
	}

	return max
}

// Sum of values
func Sum[N Number](n ...N) N {
	var sum N
	for i := range n {
		sum += n[i]
	}

	return sum
}

// SumFrom of iterated
func SumFrom[N Number](n Next[N]) N {
	var sum N
	for n := range From(n) {
		sum += n
	}

	return sum
}

// Lazy return a function that wraps around 't'
func Lazy[T any](t T) func() T {
	return func() T {
		return t
	}
}

// Id return t as is, without any modifications
func Id[T any](t T) T {
	return t
}

// Async wrapper around WaitGroup to run and wait for the completion of all goroutines
func Async() Pairs[Run, Wait] {
	var rw sync.WaitGroup
	var run = func(fn func()) {
		rw.Add(1)
		go func() {
			fn()
			rw.Done()
		}()
	}

	return Pair(run, rw.Wait)
}

// Async wrapper around WaitGroup to run and wait for the completion of all goroutines
// will pass a context to each goroutine.
func AsyncCtx(ctx context.Context) Pairs[RunCtx, Wait] {
	var rw sync.WaitGroup
	var runctx = func(fn func(ctx context.Context)) {
		rw.Add(1)
		go func() {
			fn(ctx)
			rw.Done()
		}()
	}

	return Pair(runctx, rw.Wait)
}

// DoNothing maybe useful in some cases, does nothing.
func DoNothing() {}
