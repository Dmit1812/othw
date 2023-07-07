package hw04lrucache

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	mu    sync.Mutex
	front *ListItem
	back  *ListItem
}

// Len returns the length of the list.
func (l *list) Len() int {
	// critical section start, it ends when return is called
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.len
}

// Front returns address of the front ListItem of the list.
func (l *list) Front() *ListItem {
	// critical section start, it ends when return is called
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.front
}

// Back returns address of the back ListItem of the list.
func (l *list) Back() *ListItem {
	// critical section start, it ends when return is called
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.back
}

// pushFront private method adds an already created ListItem to the front of the list
// BEWARE! it assumes that the mutex lock should be set by the caller
func (l *list) pushFront(i *ListItem) *ListItem {
	// exit if i is nil
	if i == nil {
		return nil
	}

	// initialize the new front of the list
	i.Prev = nil
	i.Next = l.front

	// if our front is initialized with pointer set old front to point to i (new front)
	if l.front != nil {
		l.front.Prev = i
	} else {
		// if our front is nil then it's the first element and both front and back should point to it
		l.back = i
	}

	// set front to the i
	l.front = i

	// increase array size
	l.len++
	return i
}

// PushFront adds a new value to the back of the list.
func (l *list) PushFront(v interface{}) *ListItem {
	// critical section start, it ends when return is called
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.pushFront(&ListItem{Value: v, Next: nil, Prev: nil})
}

// PushBack adds a new value to the back of the list.
func (l *list) PushBack(v interface{}) *ListItem {
	// critical section start
	l.mu.Lock()
	defer l.mu.Unlock()

	// create a new ListItem with value from v
	newback := &ListItem{Value: v, Next: nil, Prev: l.back}

	// if out back is initilized then make sure that previous element points to new back element
	if l.back != nil {
		l.back.Next = newback
	} else {
		// if our back is nil then newback is the first element and both front and back should point to it
		l.front = newback
	}
	// set back to the newback element
	l.back = newback
	// increase array size
	l.len++
	return newback
}

// remove removes a ListItem from the list
// no checking if the item is in the list performed
// BEWARE! it assumes that the mutex lock should be set by the caller
func (l *list) remove(i *ListItem) {
	// exit if i is nil
	if i == nil {
		return
	}

	// if previous element of element i is defined, then make sure that it points to the element that is after alement i
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		// there were no previous element, so we are removing the first element. Next element to current becomes new first
		l.front = i.Next
	}
	// if next element of element i is defined, then make sure that it points back to the element that is before i
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		// there was no next element, so we are removing the last element. Previous element to current becomes new last
		l.back = i.Prev
	}
	// we removed an element from the list reduce array size
	l.len--
}

// Remove will set mutex lock and call the remove to remove a ListItem from the list
func (l *list) Remove(i *ListItem) {
	// critical section start, it ends when return is called
	l.mu.Lock()
	defer l.mu.Unlock()

	l.remove(i)
}

// MoveToFront moves a ListItem to the front of the list
// no checking if the item is from the the list performed
func (l *list) MoveToFront(i *ListItem) {
	// critical section start, it ends when return is called
	l.mu.Lock()
	defer l.mu.Unlock()

	l.remove(i)
	l.pushFront(i)
}

// NewList creates a new list
func NewList() List {
	return new(list)
}
