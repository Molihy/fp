package slice

import (
	"fmt"
	"testing"

	"github.com/molikatty/fp"
)

func TestExample(t *testing.T) {
	t.Run("Filter", func(t *testing.T) {
		var ints = Iter(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
		t.Log(fp.Slice(fp.Filter[int](ints, func(i int) bool {
			return i > 5
		})))
	})

	t.Run("Map", func(t *testing.T) {
		var ints = Iter(Of(1, 2, 3, 4, 6, 7, 8, 9, 10))
		t.Log(fp.Slice(fp.Map[string](ints, func(i int) string {
			return fmt.Sprint(i)
		})))
	})

	t.Run("Zip", func(t *testing.T) {
		var i1 = Iter(Of(1, 2, 3, 4))
		var i2 = Iter(Of(5, 6, 7, 8))
		var i3 = Iter(Of(9, 10, 11))
		t.Log(fp.Slice(fp.Zip(i1, i2, i3)))
	})

	t.Run("Reduce", func(t *testing.T) {
		var ints = Iter(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
		t.Log(fp.Reduce(ints, func(n, m int) int {
			return n + m
		}))
	})

	t.Run("Fold", func(t *testing.T) {
		var ints = Iter(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
		t.Log(fp.Slice(fp.Fold(ints, func(n, m int) int {
			return n + m
		})))
	})

	t.Run("All", func(t *testing.T) {
		var ints = Iter(Of(true, true, true, true, true, true, true, true, true))
		t.Log(fp.All(ints))
	})

	t.Run("Single", func(t *testing.T) {
		var ints = Iter(Of(false, false, false, false, true, false, false, false))
		t.Log(fp.Single(ints))
	})

	t.Run("Merge", func(t *testing.T) {
		var i1 = Iter(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11))
		var i2 = Iter(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12))
		var i3 = Iter(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 13))
		t.Log(fp.Slice(fp.Merge(i1, i2, i3)))
	})
}
