package audio

import (
	"container/list"
	"sync"
)

type Queue struct {
	list    *list.List
	current *list.Element
	size    int
	mu      sync.RWMutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		list: list.New(),
		size: size,
	}
}

func (q *Queue) Add(ad *AudioFile) *list.Element {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.list.Len() == q.size {
		q.list.Remove(q.list.Front())
	}

	e := q.list.PushBack(ad)

	if q.current == nil {
		q.current = e
	}

	return e
}

func (q *Queue) AddAfter(ad *AudioFile, e *list.Element) *list.Element {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.list.InsertAfter(ad, e)
}

func (q *Queue) Next() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.current != nil {
		next := q.current.Next()

		if next != nil {
			q.current = next
		}
	}
}

func (q *Queue) Previous() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.current != nil {
		prev := q.current.Prev()

		if prev != nil {
			q.current = prev
		}

		q.current.Value.(*AudioFile).Play()
	}
}

func (q *Queue) MoveAfter(e *list.Element, mark *list.Element) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list.MoveAfter(e, mark)
}

func (q *Queue) Play(e *list.Element) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.current == e {
		q.current.Value.(*AudioFile).Play()
	} else {
		for c := q.list.Front(); c != nil; c = c.Next() {
			if c == e {
				q.current = c
				c.Value.(*AudioFile).Play()
				break
			}
		}
	}
}

func (q *Queue) Pause() {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.current != nil {
		q.current.Value.(*AudioFile).Pause()
	}
}

func (q *Queue) Resume() {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.current != nil {
		q.current.Value.(*AudioFile).Resume()
	}
}

func (q *Queue) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list.Init()
}
