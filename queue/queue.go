/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package queue

import (
	"container/list"
	"errors"
	"sync"
)

type Queue struct {
	lock sync.Mutex

	queueCond *sync.Cond
	queue     *list.List
	maxLen    int

	tasksDone         *sync.Cond
	unfinishedTaskCnt int
}

func (q *Queue) Init() {
	q.queueCond = sync.NewCond(&q.lock)
	q.tasksDone = sync.NewCond(&q.lock)
	q.queue = list.New()
	q.maxLen = -1
}

func (q *Queue) SetMaxLen(maxLen int) {
	q.lock.Lock()
	q.maxLen = maxLen
	q.lock.Unlock()
}

func (q *Queue) Append(item interface{}) error {
	var err error

	q.queueCond.L.Lock()

	if q.maxLen != -1 && q.queue.Len() >= q.maxLen {
		err = errors.New("Queue is full")
	} else {
		q.queue.PushBack(item)
		q.unfinishedTaskCnt++
		q.queueCond.Signal()
		err = nil
	}

	q.queueCond.L.Unlock()
	return err
}

func (q *Queue) Remove() interface{} {
	q.queueCond.L.Lock()

	for q.queue.Len() == 0 {
		q.queueCond.Wait()
	}

	item := q.queue.Front()
	q.queue.Remove(item)

	q.queueCond.L.Unlock()

	return item.Value
}

func (q *Queue) Len() int {
	var len int
	q.lock.Lock()
	len = q.queue.Len()
	q.lock.Unlock()

	return len
}

func (q *Queue) TaskDone() {
	q.tasksDone.L.Lock()

	newCnt := q.unfinishedTaskCnt - 1
	if newCnt <= 0 {
		q.tasksDone.Broadcast()
	}
	q.unfinishedTaskCnt = newCnt

	q.tasksDone.L.Unlock()
}

func (q *Queue) Join() {
	q.tasksDone.L.Lock()

	for q.unfinishedTaskCnt > 0 {
		q.tasksDone.Wait()
	}

	q.tasksDone.L.Unlock()
}
