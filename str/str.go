package str

import (
	"crypto/md5"
	"strings"

	"github.com/molikatty/fp"
	"github.com/molikatty/fp/slice"
	"github.com/spaolacci/murmur3"
)

type Char interface {
	~rune | ~byte
}

// Iter create an iterator for string.
func Iter[C Char](s string) fp.Next[C] {
	var bytes = fp.If(fp.Is[byte](fp.Zero[C]()),
		func() []C {
			return fp.To[[]C](To[byte](s))
		},
		func() []C {
			return fp.To[[]C](To[rune](s))
		},
	)

	return slice.Iter(bytes)
}

// Cat efficiently concatenate multiple strings.
func Cat(s ...string) string {
	var newstr = make([]byte, 0, fp.Sum(Lens(s...)...))

	for i := range s {
		newstr = append(newstr, s[i]...)
	}

	return String[[]byte, byte](newstr)
}

func Clone(s string) string {
	return Cat(s)
}

// Join concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func Join(sep string, s ...string) string {
	var n = fp.Sum(Lens(s...)...)

	var str = func() string {
		ns := append(make([]byte, 0, n+(len(s)-1)*fp.Sum(Lens(sep)...)), s[0]...)
		for _, v := range s[1:] {
			ns = append(append(ns, sep...), v...)
		}

		return String[[]byte, byte](ns)
	}

	return fp.If(fp.IsZero(n), fp.Zero[string], str)
}

// Wrap a string with sep.
func Warp(s, sep string) string {
	return Cat(sep, s, sep)
}

// Unwarp a string with sep.
func Unwarp(s, sep string) string {
	return fp.If(
		fp.Not(strings.HasPrefix(s, sep) || strings.HasSuffix(s, sep)),
		fp.Lazy(s),
		func() string {
			return string(s[len(sep):][len(s)-len(sep)<<1:])
		},
	)

}

// IsEmpty check a string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// Hash conver a string to hash
func Hash[N uint32 | uint64](s string) N {
	var u32 = func() any {
		return fp.Any(murmur3.Sum32(To[byte](s)))
	}

	var u64 = func() any {
		return fp.Any(murmur3.Sum64(To[byte](s)))
	}

	return fp.To[N](fp.If(fp.Is[uint32](fp.Zero[N]()), u32, u64))
}

// Md5 conver a string to md5 string
func Md5(s string) string {
	var data = md5.New()
	data.Write(To[byte](s))
	return String[[]byte, byte](data.Sum(nil))
}

// Lens store the len of multiple string into a int-slice.
func Lens(s ...string) []int {
	var n = make([]int, 0, 10)
	for i := range s {
		n = append(n, len(s[i]))
	}

	return n
}

// Contains reports whether a string is within multiple strings.
func Contains(strs []string, str string) bool {
	for i := range strs {
		if !strings.Contains(strs[i], str) {
			return false
		}
	}

	return true
}

// Reverse a string
func Reverse(s string) string {
	var rn = To[rune](s)

	for f, t := 0, len(rn)-1; f < t; f, t = f+1, t-1 {
		rn[f], rn[t] = fp.Swap(rn[f], rn[t])
	}

	return String[[]rune, rune](rn)
}

// Hide a substring in a main string using the given replacement string.
func Hide(s, sep string, start, end int) string {
	st, en, le := fp.To[uint](start), fp.To[uint](end), fp.To[uint](len(s))

	return fp.If(st > le-1 || st > en || IsEmpty(sep),
		fp.Lazy(s),
		func() string {
			return Cat(s[:start], strings.Repeat(sep, end-start), s[end:])
		},
	)

}

// Reverses some strings, all is false to reverse slice only.
func Reverses(s []string, all bool) []string {
	var fn = fp.If(all, fp.Lazy(Reverse), fp.Lazy(func(s string) string {
		return s
	}))

	for f, t := 0, len(s)-1; f < t; f, t = f+1, t-1 {
		s[f], s[t] = fp.Swap(fn(s[f]), fn(s[t]))
	}

	return s
}

// To conver a string to fp.Bytes.
//
// Warning: unsafe conver don't modify original type data.
func To[C Char](s string) []C {
	var r = func() []C {
		return fp.To[[]C]([]rune(s))
	}

	var b = func() []C {
		var head = fp.To[[2]uintptr](s)
		return fp.To[[]C]([3]uintptr{head[0], head[1], head[1]})
	}

	return fp.If(fp.Is[[]rune](fp.Zero[C]()), r, b)
}

// String conver a fp.Bytes to string
//
// Warning: unsafe conver don't modify original type data
func String[S fp.Next[C] | ~[]C, C Char](r S) string {
	return fp.If(fp.Is[fp.Next[C]](fp.Zero[S]()),
		func() string {
			return strHelp(fp.To[fp.Next[C]](r))
		},
		func() string {
			return fp.To[string](r)
		},
	)
}

func strHelp[C Char](iter fp.Next[C]) string {
	var s = make([]C, 0)
	fp.ForEach(iter, func(b C) bool {
		s = append(s, b)
		return true
	})

	return fp.If(fp.Is[byte](fp.Zero[C]()),
		func() string {
			return fp.To[string](s)
		},
		func() string {
			return fp.To[string](s)
		},
	)
}
