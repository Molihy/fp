package maps

import (
	"reflect"
	"testing"

	"github.com/molikatty/fp"
	"github.com/molikatty/fp/slice"
)

func TestExample(t *testing.T) {
	t.Run("Of", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		var m1 = map[int]int{1: 1, 2: 2, 3: 3}
		if !reflect.DeepEqual(m, m1) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Make", func(t *testing.T) {
		var m = Make[int, int](10)

		if !reflect.DeepEqual(len(m), len(make(map[int]int, 10))) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Clear", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		Clear(m)

		if !reflect.DeepEqual(m, map[int]int{}) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Iter", func(t *testing.T) {
		var m = Iter(Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)))
		for v := range fp.From(m) {
			t.Log(v)
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Log(IsEmpty(Of[int, int]()))
	})

	t.Run("Clone", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))

		if !reflect.DeepEqual(m, Clone(m)) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Pop", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))

		t.Log(Pop(m))
		t.Log(m)
	})

	t.Run("Copy", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		var m1 = Of[int, int]()
		Copy(m, m1)

		if !reflect.DeepEqual(m, Clone(m)) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Keys", func(t *testing.T) {
		t.Log(Keys(Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))))
	})

	t.Run("Values", func(t *testing.T) {
		t.Log(Values(Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))))
	})

	t.Run("Flat", func(t *testing.T) {
		var m = Flat(
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
		)
		t.Log(m)
	})

	t.Run("Equal", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		var m1 = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		if !reflect.DeepEqual(Equal(m1, m), true) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("EqualFunc", func(t *testing.T) {
		var m = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		var m1 = Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3))
		var eq = func(v1, v2 int) bool {
			return v1 == v2
		}

		if !reflect.DeepEqual(EqualFunc(m1, m, eq), true) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Lens", func(t *testing.T) {
		var m = slice.Of(
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
			Of(fp.Pair(1, 1), fp.Pair(2, 2), fp.Pair(3, 3)),
		)

		if !reflect.DeepEqual(Lens(m...), slice.Of(3, 3, 3, 3)) {
			t.Log("Not equal")
		}

		t.Log(true)
	})
}
