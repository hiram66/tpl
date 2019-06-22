package tpl

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type TaskStatus int

const (
	StatusTaskNotStarted TaskStatus = iota
	StatusTaskRunning
	StatusTaskFinished
)

type Task interface {
	Run()

	Wait()

	Result() []interface{}

	NotifyWhenDone() (chan struct{}, error)

	Status() TaskStatus

	Tag() string
}

type dTask struct {
	job func()

	result []interface{}

	tag string

	status TaskStatus

	wg sync.WaitGroup

	notifyChan chan struct{}
}

func (t *dTask) Run() {
	t.status = StatusTaskRunning
	t.wg.Add(1)
	go func() {
		t.job()
		t.status = StatusTaskFinished
		t.wg.Done()
		if t.notifyChan != nil {
			t.notifyChan <- struct{}{}
		}
	}()
}

func (t *dTask) Wait() {
	t.wg.Wait()
}

func (t *dTask) Result() []interface{} {
	t.Wait()
	return t.result
}

func (t *dTask) Status() TaskStatus {
	return t.status
}

func (t *dTask) NotifyWhenDone() (chan struct{}, error) {

	if t.status == StatusTaskFinished  {
		return nil, errors.New("cannot set notifier channel for finished task")
	}
	t.notifyChan = make(chan struct{})
	return t.notifyChan, nil
}

func (t *dTask) Tag() string {
	return t.tag
}

func NewTask(job func()) Task {
	return &dTask{job: job, status: StatusTaskNotStarted, notifyChan:nil}
}

func getValuesSlice(args []interface{}) []reflect.Value {
	result := make([]reflect.Value, 0)

	for _, v := range args {
		result = append(result, reflect.ValueOf(v))
	}

	return result
}

func TaskFrom(tag string, fn interface{}, args ...interface{}) Task {
	task := new(dTask)

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic(fmt.Errorf("%v is not a function", fn))
	}

	task.job = func() {
		for _, v := range reflect.ValueOf(fn).Call(getValuesSlice(args)) {
			task.result = append(task.result, v.Interface())
		}
	}

	task.tag = tag
	return task
}
