package common_test

import (
	"maze/common"
	"testing"
)

func TestTaskManagerPushAtomicSuccess(t *testing.T) {
	tq := common.NewBasicTaskManager()
	tq.Push(common.NewTimePriorityTask())
	if tq.Len() != 1 {
		t.Errorf("Insert one task into queue, expect queue size to be 1\n, current length is %d", tq.Len())
	}
}

func TestTaskManagerPushMaintainOrder(t *testing.T) {
	tq := common.NewBasicTaskManager()
	t1 := common.NewTimePriorityTask()
	t2 := common.NewTimePriorityTask()
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
