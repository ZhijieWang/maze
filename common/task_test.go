package common

import (
	"testing"
)

func TestTaskQueuePushAtomicSuccess(t *testing.T) {
	tq := NewTaskQueue()
	tq.Push(TimePriorityTask{})
	if tq.Len() != 1 {
		t.Errorf("Insert one task into queue, expect queue size to be 1\n")
	}
}

func TestTaskQueuePushMaintainOrder(t *testing.T) {
	tq := NewTaskQueue()
	t1 := NewTimePriorityTask()
	t2 := NewTimePriorityTask()
	tq.Push(t2)
	tq.Push(t1)
	if t1 == t2 {
		t.Error("Input should be different\n")
	}
	if t1 != tq.Pop() {
		t.Errorf("Expect the task queue to maintain time order for out of order push\n")
	}
	if t2 != tq.Pop() {
		t.Errorf("Expect the task queue to maintain time order for out of order push\n")
	}

}
