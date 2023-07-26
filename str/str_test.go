package str

import (
	"math/rand"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) []string {
	var str []string
	for i := 0; i < n; i++ {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		str = append(str, string(b))
	}
	return str
}

func BenchmarkExample(b *testing.B) {
	var s = randomString(10)

	b.Run("Hide", func(b *testing.B) {
		var st = Cat(s...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Hide(st, "*", 3, 5)
		}
	})

	b.Run("Rune", func(b *testing.B) {
		var st = Cat(s...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			To[rune](st)
		}
	})

	b.Run("byte", func(b *testing.B) {
		var st = Cat(s...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			To[byte](st)
		}
	})

	b.Run("Cat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Cat(s...)
		}
	})

	b.Run("Join", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Join("join", s...)
		}
	})

	b.Run("Hash64", func(b *testing.B) {
		var str = Cat(s...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Hash[uint64](str)
		}
	})

	b.Run("Hash32", func(b *testing.B) {
		var str = Cat(s...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Hash[uint32](str)
		}
	})

	b.Run("Md5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Md5(s[0])
		}
	})

	b.Run("Len", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Lens(s...)
		}
	})

	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Contains(s, "abc")
		}
	})

	b.Run("Reverse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Reverse(s[0])
		}
	})

	b.Run("Reverses_false", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Reverses(s, false)
		}
	})

	b.Run("Reverses_true", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Reverses(s, true)
		}
	})
}
