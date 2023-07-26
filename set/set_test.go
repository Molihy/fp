package set

import (
	"sync"
	"testing"
)

func TestExample(t *testing.T) {
	t.Run("Of", func(t *testing.T) {
		t.Log(Of[Unsafe, int]())
		t.Log(Of[Safe, int]())
	})

	t.Run("Add", func(t *testing.T) {
		var set = Of[Safe, int]()
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(i int) {
				set.Add(i)
				wg.Done()
			}(i)
		}

		wg.Wait()

		t.Log(set)
	})

	t.Run("Difference", func(t *testing.T) {
		var d = Of[Unsafe](1, 2)
		var e = Of[Unsafe](1, 2, 3)

		t.Log(d.Difference(e))
	})

	t.Run("Intersect", func(t *testing.T) {
		var d = Of[Unsafe](1, 2)
		var e = Of[Unsafe](1, 2, 3)

		t.Log(d.Intersect(e))
	})

	t.Run("IsSubset", func(t *testing.T) {
		var d = Of[Unsafe](1, 4)
		var e = Of[Unsafe](1, 2, 3)

		t.Log(d.IsSubset(e))
		t.Log(e.IsSubset(d))
	})

	t.Run("Union", func(t *testing.T) {
		var d = Of[Unsafe](1, 4)
		var e = Of[Unsafe](1, 2, 3)

		t.Log(d.Union(e))
	})
}
