package utils

import (
	"fmt"
	"sync"
)

type StringSet map[string]struct{}

func NewStringSetFromStr(str ...string) StringSet {
	set := NewStringSet()
	set.AddAll(str...)
	return set
}

func NewStringSet() StringSet {
	return make(map[string]struct{})
}

func (set StringSet) Add(str string) {
	set[str] = struct{}{}
}

func (set StringSet) AddAll(str ...string) {
	for _, v := range str {
		set[v] = struct{}{}
	}
}

func (set StringSet) Remove(str string) {
	delete(set, str)
}

func (set StringSet) RemoveAll(str ...string) {
	for _, v := range str {
		delete(set, v)
	}
}

func (set StringSet) ToSlice() (res []string) {
	for k := range set {
		res = append(res, k)
	}
	return
}

// =============================================================================
// =============================================================================
// =============================================================================
// Queue

var emptyError = fmt.Errorf("the queue is empty")

type QueueString struct {
	sync.Mutex
	elements []string
}

func (q *QueueString) Push(e string) {
	q.Lock()
	defer q.Unlock()

	q.elements = append(q.elements, e)
}

func (q *QueueString) Pop() (string, error) {
	q.Lock()
	defer q.Unlock()

	if len(q.elements) == 0 {
		return "", emptyError
	}

	e := q.elements[0]
	q.elements = q.elements[1:]
	return e, nil
}

func (q *QueueString) Size() int {
	return len(q.elements)
}
