package tpl

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestRunAll(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) string {
		return fmt.Sprintf("%s_%s", s,s)
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()
	group.WaitAll()

	for i := 0; i <= 10; i++ {
		result, e := group.GetTaskResult(fmt.Sprintf("task_%d", i))

		if e != nil {
			log.Fatal(e)
		}

		if len(result) != 1 {
			log.Fatalf("expected result count : 1, but got : %d", len(result))
		}

		switch result[0].(type) {
		case string:
			rString := result[0].(string)
			if rString != fmt.Sprintf("%d_%d", i, i) {
				log.Fatalf("Task_%d result should be %d_%d but fount %s",i,i,i,rString)
			}
			break
		default:
			log.Fatalf("Wrong Task Result Type %v", reflect.TypeOf(result))
		}
	}
}

func TestTaskGroup_WaitForTask(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) string {
		return fmt.Sprintf("%s_%s", s,s)
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(11)

	tag := fmt.Sprintf("task_%d",n)
	e := group.WaitForTask(tag)

	if e != nil {
		log.Fatal(e)
	}

	result, e := group.GetTaskResult(tag)

	if e != nil {
		log.Fatal(e)
	}

	if len(result) != 1 {
		log.Fatalf("expected result count : 1, but got : %d", len(result))
	}

	switch result[0].(type) {
	case string:
		rString := result[0].(string)
		if rString != fmt.Sprintf("%d_%d", n, n) {
			log.Fatalf("Task_%d result should be %d_%d but fount %s",n,n,n,rString)
		}
		break
	default:
		log.Fatalf("Wrong Task Result Type %v", reflect.TypeOf(result))
	}
}


func TestTaskGroup_GetTaskResult_WhenTaskIsNotCompleted(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) string {
		time.Sleep(time.Second * 10)
		return fmt.Sprintf("%s_%s", s,s)
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(11)

	tag := fmt.Sprintf("task_%d",n)
	result, e := group.GetTaskResult(tag)

	if e != nil {
		log.Fatal(e)
	}

	if result != nil {
		log.Fatal("running task result should be nil")
	}
}

func TestTaskGroup_GetTaskResult_WhenTaskDoesNotExist(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) string {
		return fmt.Sprintf("%s_%s", s,s)
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(11) + 11
	tag := fmt.Sprintf("task_%d",n)
	_, e := group.GetTaskResult(tag)

	if e == nil {
		log.Fatal("Expected : error != nil")
	}
}

func TestTaskGroup_GetResults_WhenTasksHaveAtLeastOneResult(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) string {
		return fmt.Sprintf("%s_%s", s,s)
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()
	group.WaitAll()

	results := group.GetResults()
	for i := 0; i <= 10; i++ {
		tag := fmt.Sprintf("task_%d", i)
		result, ok := results[tag]

		if !ok {
			log.Fatalf("task %s does not exist",tag)
		}

		if result == nil || len(result) != 1 {
			log.Fatal("task should have exactly one result")
		}
	}
}

func TestTaskGroup_GetResults_WhenTasksHaveNoResult(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()
	group.WaitAll()

	results := group.GetResults()
	for i := 0; i <= 10; i++ {
		tag := fmt.Sprintf("task_%d", i)
		result, ok := results[tag]

		if !ok {
			log.Fatalf("task %s does not exist",tag)
		}

		if len(result) != 0 || result != nil {
			log.Fatal("task should have exactly zero result")
		}
	}
}

func TestTaskGroup_GetTaskState(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()
	group.WaitAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	tag := fmt.Sprintf("task_%d",n)
	status, e := group.GetTaskState(tag)

	if e != nil {
		log.Fatal(e)
	}

	if status != StatusTaskFinished {
		log.Fatal("task status should be finished")
	}
}

func TestTaskGroup_GetTaskState_WhenTaskDoesNotExist(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()
	group.WaitAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(11) + 11
	tag := fmt.Sprintf("task_%d",n)
	status, e := group.GetTaskState(tag)

	if status != -1 || e == nil {
		log.Fatal("task status should be -1 and error value should not be nil")
	}
}

func TestTaskGroup_GetState(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()
	group.WaitAll()

	if !group.GetState() {
		log.Fatal("all tasks should be finished")
	}
}

func TestTaskGroup_GetState_WhenTasksAreNotFinished(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(10) + 1))
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()

	if group.GetState() {
		log.Fatal("some tasks should not be finished")
	}
}

// TODO * 2
func TestTaskGroup_RunTask(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(10) + 1))
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)


	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	tag := fmt.Sprintf("task_%d",n)
	e := group.RunTask(tag)

	if e != nil {
		log.Fatal(e)
	}

	status, e := group.GetTaskState(tag)

	if e != nil {
		log.Fatal(e)
	}

	if status != StatusTaskRunning {
		log.Fatal("task should be on running state")
	}
}

func TestTaskGroup_RunTask_WhenTaskDoesNotExist(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(10) + 1))
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(11) + 11
	tag := fmt.Sprintf("task_%d",n)
	e := group.RunTask(tag)

	if e == nil {
		log.Fatal("error value should be nil")
	}
}

func TestTaskGroup_WaitForTask_WhenTaskDoesNotExists(t *testing.T) {
	tasks := make([]Task, 0)
	fn := func(s string) string {
		return fmt.Sprintf("%s_%s", s,s)
	}
	for i := 0; i <= 10; i++ {
		t := TaskFrom(fmt.Sprintf("task_%d",i),fn, fmt.Sprintf("%d",i))
		tasks = append(tasks, t)
	}

	group := NewTaskGroup(tasks)
	group.RunAll()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(11) + 11

	tag := fmt.Sprintf("task_%d",n)
	e := group.WaitForTask(tag)

	if e == nil {
		log.Fatal("error value should not be nil")
	}
}