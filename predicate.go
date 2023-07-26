package fp

// If like in Python, the expression "True and 1 or -1" is used to achieve
// the effect of short-circuit ternary operation. It works by evaluating
// the condition, and if the condition is true, it returns the value after
// "and" (1 in this case). If the condition is false, it returns the value
// after "or" (-1 in this case).
func If[T any](ok bool, t, f func() T) T {
	if ok {
		return t()
	}

	return f()
}

// OrElse returns the value of b if 'a' is zero, otherwise returns a”
func OrElse[T comparable](a T, b func() T) T {
	return If(IsZero(a), b, Lazy(a))
}

// Def check ok, if it is true, return the value of b, otherwise return the zero value.
func Def[T any](ok bool, b func() T) T {
	return If(ok, b, Zero[T])
}

// IsNil checks if the value 't' is nil.
func IsNil[T any](t T) bool {
	return To[uintptr](t) == 0
}

// NotNil checks if the value 't' not nil.
func NotNil[T any](t T) bool {
	return Not(IsNil(t))
}

func IsNaN[T Size](t T) bool {
	return t != t
}

func NotNaN[T Size](t T) bool {
	return Not(IsNaN(t))
}

// IsZero checks if the value 't' is default zero value.
func IsZero[T comparable](t T) bool {
	return t == Zero[T]()
}

// IsZero checks if the value 't' is not default zero value.
func NotZero[T comparable](n T) bool {
	return Not(IsZero(n))
}

// True return a true
func True() bool {
	return true
}

// False return a false
func False() bool {
	return false
}

// Is assert that an empty interface is of a specific type
func Is[T any](t any) bool {
	_, ok := t.(T)
	return ok
}

// Has check if a value exists in a map
func Has[T comparable, R any](m map[T]R, t T) bool {
	_, ok := m[t]
	return ok
}

// Or if function f returns true, then function b is not executed,
// which is an implementation of the '||' operator in a function form
func Or(f, s func() bool) func() bool {
	return func() bool {
		return f() || s()
	}
}

// And if function f returns true, then function b is executed,
// which is an implementation of the '&&' operator in a function form.
func And(a, b func() bool) func() bool {
	return func() bool {
		return a() && b()
	}
}

// Xor performs the XOR operation on two bool values.
func Xor(a, b bool) bool {
	return To[bool](To[uint8](a) ^ To[uint8](b))
}

// Not against a bool
func Not(b bool) bool {
	return !b
}

// All check if there is any 'false' in the bool iterator, similar to Python
// built-in function 'all'.
func All(iter Next[bool]) (ok bool) {
	ForEach(iter, func(b bool) bool {
		ok = b
		return b
	})

	return
}
