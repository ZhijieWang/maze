package common

import (
	"testing"
)

func TestCanModifyRobot(t *testing.T) {
	g := CreateWorld(3)
	r := g.GetRobots()
	r[0].location = r[1].location
	if r[0].location != g.GetRobots()[0].location {
		t.Errorf("Expect the robots returned to be modifiable\n")
	}
}

func TestCanModifyTasks(t *testing.T) {
	g := CreateWorld(3)

	task := make([]*Task, 1)
	task = append(task, &Task{})
	g.SetTasks(task)
	if len(g.GetTasks()) == 0 {
		t.Errorf("Expect the task list to be mutable\n")
	}
	if g.GetTasks()[0] != task[0] {
		t.Errorf("Expect the Task setter methog to work, but failed")
	}
}
