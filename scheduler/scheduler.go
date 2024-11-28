package scheduler

import (
	"container/heap"
	"time"
)

type Scheduler struct {
	tasks     tasksHeap
	idCounter uint16
}

type TaskFunc func() (completed bool, err error)

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) ScheduleOnce(
	delay time.Duration,
	tf TaskFunc,
) uint16 {
	s.idCounter++

	s.tasks = append(s.tasks, task{
		id:   s.idCounter,
		fn:   tf,
		date: time.Now().Add(delay),
	})

	return s.idCounter
}

func (s *Scheduler) ScheduleRepeating(
	delay time.Duration,
	tf TaskFunc,
) uint16 {
	s.idCounter++

	heap.Push(&s.tasks, task{
		id:     s.idCounter,
		fn:     tf,
		date:   time.Now().Add(delay),
		repeat: delay,
	})

	return s.idCounter
}

func (s *Scheduler) Cancel(id uint16) {
	for i, t := range s.tasks {
		if t.id == id {
			heap.Remove(&s.tasks, i)

			return
		}
	}
}

func (s *Scheduler) Run() error {
	if len(s.tasks) == 0 || s.tasks[0].date.After(time.Now()) {
		return nil
	}

	for {
		item := heap.Pop(&s.tasks)
		if item == nil {
			return nil
		}

		t := item.(task)

		if t.date.After(time.Now()) {
			heap.Push(&s.tasks, t)

			return nil
		}

		completed, err := t.fn()
		if err != nil {
			return err
		}
		if !completed && t.repeat > 0 {
			heap.Push(&s.tasks, task{
				id:     t.id,
				fn:     t.fn,
				date:   time.Now().Add(t.repeat),
				repeat: t.repeat,
			})
		}
	}
}

type task struct {
	id     uint16
	fn     TaskFunc
	date   time.Time
	repeat time.Duration
}

type tasksHeap []task

func (h tasksHeap) Len() int {
	return len(h)
}

func (h *tasksHeap) Less(i, j int) bool {
	return (*h)[i].date.Before((*h)[j].date)
}

func (h *tasksHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *tasksHeap) Push(x any) {
	*h = append(*h, x.(task))
}

func (h *tasksHeap) Pop() any {
	old := *h

	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
