package heap_test

import (
	"math/rand"
	"slices"
	"strconv"
	"testing"

	"github.com/rinnothing/pinkerton/pkg/heap"
)

func compare(a, b int) bool {
	return a < b
}

func getHeap() heap.Heap[int, string] {
	return heap.New[int, string](compare)
}

func addXRandom(h heap.Heap[int, string], arr *[]int, n int) {
	for range n {
		val := rand.Int()
		h.Push(strconv.Itoa(val), val)
		*arr = append(*arr, val)
	}
	slices.Sort(*arr)
}

func retrieveXAndCheckOrder(t *testing.T, h heap.Heap[int, string], arr *[]int, n int) {
	for range n {
		fstVal, fstKey, ok := h.Pop()
		if !ok {
			t.Fatal("heap shouldn't be empty")
		}

		if strconv.Itoa(fstVal) != fstKey {
			t.Errorf("val and key should be equal but they're not: %d != %s", fstVal, fstKey)
		}

		if (*arr)[0] != fstVal {
			t.Fatalf("order is broken first item should be %d, but is %d", fstVal, (*arr)[0])
		}
		*arr = (*arr)[1:]
	}
}

func TestOrder(t *testing.T) {
	beg, batch, times := 10, 5, 20

	var arr []int
	h := getHeap()
	addXRandom(h, &arr, beg)

	for range times {
		retrieveXAndCheckOrder(t, h, &arr, batch)
		addXRandom(h, &arr, batch)
	}

	retrieveXAndCheckOrder(t, h, &arr, beg)
	_, _, ok := h.Pop()
	if ok {
		t.Fatalf("heap should be empty by now, but it's not, size = %d", h.Len()+1)
	}
}

func TestPopTopEquality(t *testing.T) {
	sz := 10

	var arr []int
	h := getHeap()
	addXRandom(h, &arr, sz)

	for range sz + 1 {
		topVal, topKey, topOk := h.Top()
		popVal, popKey, popOk := h.Pop()

		if topVal != popVal {
			t.Fatalf("top and pop vals should be equal, but they're not: %d != %d", topVal, popVal)
		}

		if topKey != popKey {
			t.Fatalf("top and pop keys should be equal, but they're not: %s != %s", topKey, popKey)
		}

		if topOk != popOk {
			t.Fatalf("top and pop keys should be equal, but they're not: %t != %t", topOk, popOk)
		}
	}
}

func TestGet(t *testing.T) {
	sz := 10
	var arr []int
	h := getHeap()
	addXRandom(h, &arr, sz)

	for range sz {
		val := arr[rand.Int()%len(arr)]
		key := strconv.Itoa(val)
		res, ok := h.Get(key)
		if !ok {
			t.Fatalf("haven't found key %s in heap", key)
		}
		if res != val {
			t.Fatalf("value of key %s isn't correct, should be %d, but is %d", key, val, res)
		}

		_, _, _ = h.Pop()
		arr = arr[1:]
	}
}

func TestRemove(t *testing.T) {
	sz := 10
	var arr []int
	h := getHeap()
	addXRandom(h, &arr, sz)

	for range sz {
		idx := rand.Int() % len(arr)
		val := arr[idx]
		key := strconv.Itoa(val)
		h.Remove(key)
		_, ok := h.Get(key)
		if ok {
			t.Fatalf("have found key %s in heap, even though it was deleted", key)
		}

		h.Pop()

		for i := idx; i < len(arr)-1; i++ {
			arr[i] = arr[i+1]
		}
		arr = arr[:len(arr)-1]
	}
}

func TestLen(t *testing.T) {
	sz := 10
	var arr []int
	h := getHeap()
	addXRandom(h, &arr, sz)

	for s := range sz {
		szNow := sz - s
		hLen := h.Len()
		if szNow != hLen {
			t.Fatalf("heap has wrong size, should be %d, but is %d", szNow, hLen)
		}

		h.Pop()
	}
}
