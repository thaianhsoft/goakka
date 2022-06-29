package goakka

import "fmt"

type Element struct {
	data interface{}
	next *Element
}
type Queue struct {
	front *Element
	back *Element
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) PushBack(data interface{}) {
	e := Element{data: data}
	if q.front == nil {
		q.front, q.back = &e, &e
	} else {
		old := q.back
		old.next, q.back = &e, &e
	}
}

func (q *Queue) PopFront() interface{} {
	popNode := q.front
	q.front = popNode.next
	return popNode.data
}

func (q *Queue) Empty() bool {
	return q.front == nil
}

func (q *Queue) PushFront(data interface{}) {
	e := Element{data:data}
	if q.front == nil {
		q.front = &e
		q.back = &e
	} else {
		fmt.Println("push here")
		old_front := q.front
		q.front, e.next = &e, old_front
	}
}

