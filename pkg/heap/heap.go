package heap

import "container/heap"

type Heap[T any, K comparable] interface {
	Push(key K, val T)
	Pop() (T, K, bool)
	Top() (T, K, bool)
	Get(key K) (T, bool)
	Remove(key K)
	Len() int
}

var _ Heap[any, any] = &heapImpl[any, any]{}

func New[T any, K comparable](compare func(a, b T) bool) *heapImpl[T, K] {
	return &heapImpl[T, K]{container: &heapContainer[T, K]{compare: compare, arr: make([]*pair[K, T], 0), keyToPlaces: make(map[K]int), len: 0}}
}

type pair[K, V any] struct {
	key K
	val V
}

type heapImpl[T any, K comparable] struct {
	container *heapContainer[T, K]
}

func (h *heapImpl[T, K]) Len() int {
	return h.container.Len()
}

func (h *heapImpl[T, K]) Pop() (T, K, bool) {
	if h.container.len == 0 {
		var defVal T
		var defKey K
		return defVal, defKey, false
	}

	val := heap.Pop(h.container)
	res, ok := val.(*pair[K, T])
	if !ok {
		panic("can't cast value to it's type")
	}

	return res.val, res.key, true
}

func (h *heapImpl[T, K]) Top() (T, K, bool) {
	if h.container.len == 0 {
		var defVal T
		var defKey K
		return defVal, defKey, false
	}

	val := h.container.arr[0]
	return val.val, val.key, true
}

func (h *heapImpl[T, K]) Get(key K) (T, bool) {
	place, ok := h.container.keyToPlaces[key]
	if !ok {
		var def T
		return def, false
	}

	res := h.container.arr[place]
	return res.val, true
}

func (h *heapImpl[T, K]) Push(key K, val T) {
	place, ok := h.container.keyToPlaces[key]
	if !ok {
		heap.Push(h.container, &pair[K, T]{key: key, val: val})
		return
	}

	h.container.arr[place] = &pair[K, T]{key: key, val: val}
	heap.Fix(h.container, place)
}

func (h *heapImpl[T, K]) Remove(key K) {
	place, ok := h.container.keyToPlaces[key]
	if !ok {
		return
	}

	heap.Remove(h.container, place)
}

var _ heap.Interface = &heapContainer[any, any]{}

type heapContainer[T any, K comparable] struct {
	compare     func(a, b T) bool
	arr         []*pair[K, T]
	keyToPlaces map[K]int
	len         int
}

func (h *heapContainer[T, K]) Len() int {
	return h.len
}

func (h *heapContainer[T, K]) Less(i int, j int) bool {
	return h.compare(h.arr[i].val, h.arr[j].val)
}

func (h *heapContainer[T, K]) Pop() any {
	h.len--
	val := h.arr[h.len]
	h.arr = h.arr[:h.len]

	delete(h.keyToPlaces, val.key)

	return val
}

func (h *heapContainer[T, K]) Push(x any) {
	val, ok := x.(*pair[K, T])
	if !ok {
		panic("can't cast value to it's type")
	}

	h.arr = append(h.arr, val)
	h.keyToPlaces[val.key] = h.len
	h.len++
}

func (h *heapContainer[T, K]) Swap(i int, j int) {
	iVal, jVal := h.arr[i], h.arr[j]
	h.keyToPlaces[iVal.key] = j
	h.keyToPlaces[jVal.key] = i
	h.arr[j], h.arr[i] = iVal, jVal
}
