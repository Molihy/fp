package channel

import "github.com/molikatty/fp"

// Lazy initialize a channel
// When the terminating channel ends, it will be automatically closed.
func Lazy[C ~chan T, T any](c C, fn func(C)) C {
	go func() {
		fn(c)
		close(c)
	}()

	return c
}

func Make[C ~chan T, T any](bufcap int) C {
	return make(C, bufcap)
}

// Of quick make a channel
// When the terminating channel ends, it will be automatically closed.
func Of[T any](t ...T) chan T {
	var lazy = func(ch chan T) {
		for i := range t {
			ch <- t[i]
		}
	}

	return Lazy(make(chan T, len(t)), lazy)
}

// Iter make an iterator for channel
func Iter[C ~chan T, T any](c C) fp.Next[T] {
	return func() (T, bool) {
		for t := range c {
			return t, true
		}

		return fp.Zero[T](), false
	}
}

func From[C ~chan T, T any](next fp.Next[T]) C {
	return fp.Chan(next, 0)
}

// Merge multiple channels into one channel
func Merge[C ~chan T, T any](ups ...C) C {
	var lazy = func(s C) {
		r, w := fp.Async().Expand()
		output := func(up C) {
			for v := range up {
				s <- v
			}
		}

		for i := range ups {
			index := i
			r(func() { output(ups[index]) })
		}

		w()
	}

	return Lazy(make(C, fp.Sum(Caps(ups...)...)), lazy)
}

// Split the output of a channel into two
func Split[C ~chan T, T any](in C) fp.Pairs[C, C] {
	var o1, o2 = make(C), make(C)

	go func() {
		for v := range in {
			select {
			case o1 <- v:
			case o2 <- v:
			}
		}
		close(o1)
		close(o2)
	}()

	return fp.Pair(o1, o2)
}

// Caps store the cap of multiple channel into a int-slice.
func Caps[C ~chan T, T any](pipes ...C) []int {
	var n = make([]int, len(pipes))
	for i := range pipes {
		n[i] = cap(pipes[i])
	}

	return n
}

// Lens store the len of multiple channel into a int-slice.
func Lens[C ~chan T, T any](pipes ...C) []int {
	var n = make([]int, len(pipes))
	for i := range pipes {
		n[i] = len(pipes[i])
	}

	return n
}

// Buffer create a buffer for a channel
func Buffer[C ~chan T, T any](t C, bufcap int) chan T {
	var bcap = fp.If(bufcap < 0, fp.Zero[int], fp.Lazy(bufcap))

	var lazy = func(ch C) {
		for v := range t {
			ch <- v
		}
	}

	return Lazy(make(C, bcap), lazy)
}
