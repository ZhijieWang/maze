/*
 *  Copyright (c) 2019 Zhijie (Bill) Wang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package task

import (
	"errors"
	"maze/common"
	"sync"
)

type SimulatedTaskManager struct {
	tasks   map[common.TaskID]common.Task
	active  map[common.TaskID]common.Task
	archive map[common.TaskID]common.Task
}

func CreateSimulatedTaskManager() *SimulatedTaskManager {
	return &SimulatedTaskManager{
		make(map[common.TaskID]common.Task),
		make(map[common.TaskID]common.Task),
		make(map[common.TaskID]common.Task),
	}
}
func (stm *SimulatedTaskManager) GetBroadcastInfo() interface{} {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) GetAllTasks() []common.Task {
	values := []common.Task{}
	for _, t := range stm.tasks {
		values = append(values, t)
	}
	return values
}

func (stm *SimulatedTaskManager) GetNext() common.Task {
	if len(stm.tasks) == 0 {
		return nil
	}
	for _, v := range stm.tasks {
		if v.GetStatus() != common.Assigned {
			return v
		}

	}
	return nil
}

func (stm *SimulatedTaskManager) GetTasks(n int) []common.Task {
	values := make([]common.Task, n)
	for _, t := range stm.tasks {
		n -= 1
		values = append(values, t)
		if n == 0 {
			break
		}
	}
	return values
}

func (stm *SimulatedTaskManager) TaskUpdate(taskID common.TaskID, status common.TaskStatus) error {

	switch status {

	case common.Completed:
		if t, ok := stm.active[taskID]; ok {
			stm.archive[taskID] = t
			delete(stm.active, taskID)
			return nil
		} else {
			return errors.New("status can't jump from UnAssigned to Completed")
		}

	case common.Assigned:
		if t, ok := stm.tasks[taskID]; ok {
			stm.active[taskID] = t
			delete(stm.tasks, taskID)
			return nil
		} else {
			return errors.New("task not found")
		}

	default:
		return nil
	}
}

func (stm *SimulatedTaskManager) AddTask(t common.Task) bool {
	if t.GetStatus() == common.Completed {
		return false
	} else {
		if _, ok := stm.tasks[t.GetTaskID()]; ok {
			// task already in the tracker
			// edge case, return false for now
			return false
		} else {
			stm.tasks[t.GetTaskID()] = t
			return true
		}
	}
}

func (stm *SimulatedTaskManager) AddTasks(tList []common.Task) bool {
	result := true
	for _, t := range tList {
		result = result && stm.AddTask(t)
	}
	return result
}

func (stm *SimulatedTaskManager) HasTasks() bool {
	return len(stm.tasks) != 0
}
func (stm *SimulatedTaskManager) FinishedCount() int {
	return len(stm.archive)
}

func (stm *SimulatedTaskManager) ActiveCount() int {
	return len(stm.active)
}

type SimulatedTaskManagerSync struct {
	tasks   sync.Map
	active  sync.Map
	archive sync.Map
}

func (stm *SimulatedTaskManagerSync) GetBroadcastInfo() interface{} {
	panic("not implemented")
}

func (stm *SimulatedTaskManagerSync) GetAllTasks() []common.Task {
	var values []common.Task

	stm.tasks.Range(func(key interface{}, value interface{}) bool {
		values = append(values, value.(common.Task))
		return true
	})

	return values
}

func (stm *SimulatedTaskManagerSync) GetNext() common.Task {

	tasks := stm.GetAllTasks()
	if len(tasks) == 0 {
		return nil
	} else {
		return tasks[0]
	}
}

func (stm *SimulatedTaskManagerSync) GetTasks(n int) []common.Task {
	return stm.GetAllTasks()[:n]
}

func (stm *SimulatedTaskManagerSync) TaskUpdate(taskID common.TaskID, status common.TaskStatus) error {

	switch status {

	case common.Completed:
		if t, ok := stm.active.Load(taskID); ok {
			stm.archive.Store(taskID, t)
			stm.active.Delete(taskID)
			return nil
		} else {
			return errors.New("status can't jump from UnAssigned to Completed")
		}

	case common.Assigned:
		if t, ok := stm.tasks.Load(taskID); ok {
			stm.active.Store(taskID, t)
			stm.tasks.Delete(taskID)
			return nil
		} else {
			return errors.New("task not found")
		}

	default:
		return nil
	}
}

func (stm *SimulatedTaskManagerSync) AddTask(t common.Task) bool {
	if t.GetStatus() == common.Completed {
		return false
	} else {
		_, ok := stm.tasks.LoadOrStore(t.GetTaskID(), t)
		if ok {
			//loaded
			return false
		} else {
			// stored
			return true
		}

	}
}

func (stm *SimulatedTaskManagerSync) AddTasks(tList []common.Task) bool {
	result := true
	for _, t := range tList {
		result = result && stm.AddTask(t)
	}
	return result
}

func (stm *SimulatedTaskManagerSync) HasTasks() bool {
	return len(stm.GetAllTasks()) != 0
}
func (stm *SimulatedTaskManagerSync) FinishedCount() int {
	var values []common.Task

	stm.archive.Range(func(key interface{}, value interface{}) bool {
		values = append(values, value.(common.Task))
		return true
	})

	return len(values)
}

func (stm *SimulatedTaskManagerSync) ActiveCount() int {
	var values []common.Task

	stm.active.Range(func(key interface{}, value interface{}) bool {
		values = append(values, value.(common.Task))
		return true
	})

	return len(values)
}

func CreateSimulatedTaskManagerSync() *SimulatedTaskManagerSync {
	return &SimulatedTaskManagerSync{
		sync.Map{},
		sync.Map{},
		sync.Map{},
	}
}
