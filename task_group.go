package tpl

import "errors"

type TaskGroup interface {
	RunAll()

	WaitAll()

	WaitForTask(tag string) error

	RunTask(tag string) error

	GetState() bool

	GetTaskState(tag string) (TaskStatus, error)

	GetResults() map[string]interface{}

	GetTaskResult(tag string) (interface{}, error)
}

type taskGroup map[string]Task

func (tg taskGroup) RunAll() {
	for _, task := range tg {
		task.Run()
	}
}

func (tg taskGroup) WaitAll() {
	for _, task := range tg {
		task.Wait()
	}
}

func (tg taskGroup) WaitForTask(tag string) error {
	t, ok := tg[tag]
	if ok {
		t.Wait()
		return nil
	}

	return errors.New("task does not exists")
}

func (tg taskGroup) RunTask(tag string) error {
	t, ok := tg[tag]
	if ok {
		t.Run()
		return nil
	}

	return errors.New("task does not exists")
}

func (tg taskGroup) GetState() bool {
	for _, task := range tg {
		if task.Status() != StatusTaskFinished {
			return false
		}
	}
	return true
}

func (tg taskGroup) GetTaskState(tag string) (TaskStatus, error) {
	t, ok := tg[tag]
	if ok {
		return t.Status(), nil
	}
	return 1, errors.New("task does not exists")
}

func (tg taskGroup) GetResults() map[string]interface{} {
	tg.WaitAll()
	result := make(map[string]interface{})
	for tag, task := range tg {
		if len(task.Result()) > 0 {
			result[tag] = task.Result()
		}
	}

	return result
}

func (tg taskGroup) GetTaskResult(tag string) (interface{}, error) {
	t, ok := tg[tag]
	if ok {
		t.Wait()
		return t.Result(), nil
	}
	return nil, errors.New("task does not exists")
}

//NewTaskGroup Creates an implementation of task group
func NewTaskGroup(tasks []Task) TaskGroup {
	var tg taskGroup = make(map[string]Task)
	for _, v := range tasks {
		tg[v.Tag()] = v
	}

	return tg
}

