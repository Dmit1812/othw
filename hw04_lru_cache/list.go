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
	front *ListItem
	back  *ListItem
}

var mutex sync.Mutex

func (l *list) Len() int {
	mutex.Lock()
	defer mutex.Unlock()
	return l.len
}

func (l *list) Front() *ListItem {
	mutex.Lock()
	defer mutex.Unlock()
	return l.front
}

func (l *list) Back() *ListItem {
	mutex.Lock()
	defer mutex.Unlock()
	return l.back
}

func (l *list) pushFront(i *ListItem) *ListItem {
	mutex.Lock()
	defer mutex.Unlock()
	if i == nil {
		return nil
	}
	i.Prev = nil
	i.Next = l.front
	if l.front != nil {
		l.front.Prev = i
	} else {
		l.back = i
	}
	l.front = i
	l.len++
	return i
}

func (l *list) PushFront(v interface{}) *ListItem {
	return l.pushFront(&ListItem{Value: v, Next: nil, Prev: nil})
}

func (l *list) PushBack(v interface{}) *ListItem {
	mutex.Lock()
	defer mutex.Unlock()
	newback := &ListItem{Value: v, Next: nil, Prev: l.back}
	if l.back != nil {
		l.back.Next = newback
	} else {
		l.front = newback
	}
	l.back = newback
	l.len++
	return newback
}

func (l *list) Remove(i *ListItem) {
	mutex.Lock()
	defer mutex.Unlock()
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFront(i)
}

func NewList() List {
	return new(list)
}
