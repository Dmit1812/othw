package hw04lrucache

import (
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	// Confirm that list is consistent after adding to it in parallel goroutines
	t.Run("parallel", func(t *testing.T) {
		var wg sync.WaitGroup

		l := NewList()
		a := []int{40, 50, 60, 70, 80, 10, 1000, 15, 75, 30, 90, 1500}
		threads := 20

		// spawn `threads` goroutines where each would populate the list
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func(index int, array []int) {
				defer wg.Done()
				for i, v := range array {
					_ = i
					l.PushFront(v)
					l.PushBack(v)
				}
				_ = index
			}(i, a)
		}

		wg.Wait()

		var incorrectPointer bool

		// create a slice of values from list from front to back
		elemsfb := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elemsfb = append(elemsfb, i.Value.(int))
			if i.Prev != nil {
				if i.Prev.Next != i {
					incorrectPointer = true
				}
			}
			if i.Next != nil {
				if i.Next.Prev != i {
					incorrectPointer = true
				}
			}
		}

		// create a slice of values from list from back to front
		elemsbf := make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Prev {
			elemsbf = append(elemsbf, i.Value.(int))
			if i.Prev != nil {
				if i.Prev.Next != i {
					incorrectPointer = true
				}
			}
			if i.Next != nil {
				if i.Next.Prev != i {
					incorrectPointer = true
				}
			}
		}
		// reverse the slice elemsbf so it can be compared with elemsfb
		sort.SliceStable(elemsbf, func(i, j int) bool {
			return i > j
		})

		// it is expected that:
		// 1 - count of of elements from front to back and from back to front should be equal to list length
		// 2 - number of elements in list is equal to `threads` * length of initialization array
		require.Truef(t, len(elemsfb) == l.Len() && len(elemsbf) == l.Len() && 2*threads*len(a) == l.Len(), "number of items in List if traversed front to back (%d) should be equal to list length stored inside (%d) and to back to front (%d) and to threads * length (%d)", len(elemsfb), l.Len(), len(elemsbf), threads*len(a))
		// 3 - elements correctly reference each other and back
		require.Falsef(t, incorrectPointer, "elements should correctly reference each other and back")
		// 4 - values are same when we compare front to back and reversed back to front
		require.Equalf(t, elemsbf, elemsfb, "values shall be the same when we compare list created by front to back traversal and reversed list done back to front")
	})

	// Verify that list is consistent after removing everything from it
	t.Run("remove all", func(t *testing.T) {
		l := NewList()

		l.PushBack(20)     // [20]
		l.PushFront(10)    // [10, 20]
		l.PushBack(30)     // [10, 20, 30]
		l.Remove(l.Back()) // [10, 20]
		l.Remove(l.Back()) // [10]
		l.Remove(l.Back()) // []

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	// Ensure first pushback works correctly
	t.Run("first pushback", func(t *testing.T) {
		l := NewList()

		l.PushBack(20)  // [10]
		l.PushFront(10) // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())
	})

}
