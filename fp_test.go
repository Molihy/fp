package fp

import (
	"strconv"
	"testing"
)

func BenchmarkExample(b *testing.B) {
	b.Run("Compose", func(b *testing.B) {
		var t1 = func(n int) int {
			return n + n
		}

		var t2 = func(n int) int {
			return n + n
		}

		var t3 = func(n int) int {
			return n + n
		}

		for i := 0; i < b.N; i++ {
			Compose(t3, t2, t1)(10)
		}
	})

	b.Run("ComposeFAny", func(b *testing.B) {
		var t0 = ToFAny(func(n string) int {
			m, _ := strconv.Atoi(n)
			return m
		})

		var t1 = ToFAny(func(n int) string {
			return strconv.Itoa(n) + strconv.Itoa(n)
		})

		var t2 = ToFAny(func(n int) int {
			return n + n
		})

		var t3 = ToFAny(func(n int) int {
			return n + n
		})

		var add = Compose(t0, t1, t2, t3)
		for i := 0; i < b.N; i++ {
			add(1)
		}
	})

	b.Run("Curry", func(b *testing.B) {
		var t1 = func(a, b int) int {
			return a + b
		}

		for i := 0; i < b.N; i++ {
			_ = Curry(t1)(10, 20).(int)
		}
	})

	b.Run("Apply", func(b *testing.B) {
		var t = func(a, b, c, d int) int {
			return a + b + c + d
		}

		for i := 0; i < b.N; i++ {
			Apply(t, []any{1, 2, 3, 4})
		}
	})
}
