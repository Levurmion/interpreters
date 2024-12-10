package arrays

import "interpreters/utilities/graphs"

type Queue[T any] struct {
	linkedList *graphs.LinkedList[T]
}

func NewQueue[T any]() *Queue[T] {
	return new(Queue[T])
}

func (q *Queue[T]) Size() int {
	return q.linkedList.GetLength()
}

func (q *Queue[T]) Enqueue(item T) {
	q.linkedList.Unshift(item)
} 

func (q *Queue[T]) Dequeue() T {
	return q.linkedList.Pop().Item
}

func (q *Queue[T]) Top() T {
	return q.linkedList.GetTail().Item
}