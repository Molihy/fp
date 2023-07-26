package channel

import (
	"reflect"
	"testing"

	"github.com/molikatty/fp"
	"github.com/molikatty/fp/slice"
)

func TestExample(t *testing.T) {
	t.Run("Of", func(t *testing.T) {
		var c = Of(1, 2, 3)
		for i := range c {
			t.Log(i)
		}
	})

	t.Run("Lazy", func(t *testing.T) {
		var c = make(chan int)
		var ini = func(c chan int) {
			for i := 1; i < 4; i++ {
				c <- i
			}
		}

		for i := range Lazy(c, ini) {
			t.Log(i)
		}
	})

	t.Run("Iter", func(t *testing.T) {
		var c = Iter(Of(1, 2, 3, 4))
		for i := range fp.From(c) {
			t.Log(i)
		}
	})

	t.Run("Make", func(t *testing.T) {
		var c = Make[chan int, int](10)
		var c1 = make(chan int, 10)
		if !reflect.DeepEqual(cap(c), cap(c1)) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Merge", func(t *testing.T) {
		var c = Merge(Of(1, 2, 3), Of(1, 2, 3), Of(1, 2, 3))
		for i := range c {
			t.Log(i)
		}
	})

	t.Run("Split", func(t *testing.T) {
		var c = Split(Of(1, 2, 3, 4))
		for {
			select {
			case v, ok := <-c.Key():
				if !ok {
					t.Log(ok, "o1")
					return
				}

				t.Log(v, "o1")
			case v, ok := <-c.Value():
				if !ok {
					t.Log(ok, "o2")
					return
				}

				t.Log(v, "o2")
			}
		}
	})

	t.Run("Caps", func(t *testing.T) {
		var cn = Caps(Of(1, 2, 3), Of(1, 2, 3), Of(1, 2, 3))
		if !reflect.DeepEqual(cn, slice.Of(3, 3, 3)) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Lens", func(t *testing.T) {
		var cn = Lens(Of(1, 2, 3), Of(1, 2, 3), Of(1, 2, 3))
		if !reflect.DeepEqual(cn, slice.Of(0, 0, 0)) {
			t.Log("Not equal")
		}

		t.Log(true)
	})

	t.Run("Buffer", func(t *testing.T) {
		if !reflect.DeepEqual(cap(Buffer(Of(1, 2, 3, 4), 20)), 20) {
			t.Log("Not equal")
		}

		t.Log(true)
	})
}
