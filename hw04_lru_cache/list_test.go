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

		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func(index int, array []int) {
				defer wg.Done()
				for i, v := range array {
					if i%2 == 0 {
						l.PushFront(v)
					} else {
						l.PushBack(v)
					}
				}
				_ = index
			}(i, a)
		}

		wg.Wait()

		var incorrectPointer bool

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
		// reverse the list
		sort.SliceStable(elemsbf, func(i, j int) bool {
			return i > j
		})

		require.Equal(t, len(elemsfb), l.Len())
		require.Equal(t, len(elemsbf), l.Len())
		require.Equal(t, threads*len(a), l.Len())
		require.False(t, incorrectPointer)
		require.Equal(t, elemsbf, elemsfb)

	})

}
