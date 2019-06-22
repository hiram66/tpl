package tpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	result := "test string"
	b := make([]byte, 0)
	buffer := bytes.NewBuffer(b)

	task := NewTask(func() {
		buffer.WriteString(result)
	}).(*dTask)

	task.job()
	all, e := ioutil.ReadAll(buffer)

	if e != nil {
		log.Fatalln(e)
	}

	if string(all) != result {
		log.Fatalf("expected value %s, got %s", result, string(all))
	}
}

func TestTaskFrom_WhenTaskHasNoOutPut(t *testing.T) {
	var wg sync.WaitGroup
	testFunc := func(i *int64) {
		atomic.AddInt64(i, 1)
		wg.Done()
	}
	rand.Seed(time.Now().UnixNano())
	taskNum := rand.Intn(20)
	if taskNum <= 1 {
		taskNum = 10
	}
	i := new(int64)
	for index := 0; index < taskNum; index++ {
		wg.Add(1)
		t := TaskFrom(fmt.Sprintf("tag %d",index), testFunc,i).(*dTask)
		go t.job()
	}

	wg.Wait()

	if *i != int64(taskNum) {
		log.Fatalf("expected value %d, got %d", taskNum, *i)
	}
}

func TestTaskFrom_WhenTaskHasOutput(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testData := rand.Intn(20)
	
	fn := func(i int) int {
		return i * 2
	}
	
	task := TaskFrom("multiplier", fn, testData)
	task.Run()
	task.Wait()
	
	result := task.Result()

	switch result[0].(type) {
	case int:
		break
	default:
		log.Fatalf("expected type is %s",reflect.TypeOf(1).String())
		
	}

	if result[0].(int) != testData * 2 {
		log.Fatalf("expected result is %d, but found %d", testData * 2, result[0].(int))
	}

	if len(result) != 1 {
		log.Fatalf("expected output count is %d, but found %d",1,len(result))
	}
}

func TestTaskFrom_WhenNotifierExists(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testData := rand.Intn(20)

	fn := func(i int) int {
		time.Sleep(time.Second)
		return i * 2
	}

	task := TaskFrom("multiplier", fn, testData)
	task.Run()

	done, e := task.NotifyWhenDone()
	if e != nil {
		log.Fatal(e)
	}

	notified := new(bool)
	*notified = false
	go func(n *bool, c chan struct{}) {
		for {
			select {
			case <-c:
				*n = true
			}
		}
	}(notified, done)
	task.Wait()

	result := task.Result()

	switch result[0].(type) {
	case int:
		break
	default:
		log.Fatalf("expected type is %s",reflect.TypeOf(1).String())

	}

	if result[0].(int) != testData * 2 {
		log.Fatalf("expected result is %d, but found %d", testData * 2, result[0].(int))
	}

	if len(result) != 1 {
		log.Fatalf("expected output count is %d, but found %d",1,len(result))
	}

	if !*notified {
		log.Fatal("notification did not sent")
	}
}

func TestTaskFrom_WhenInputIsNotFunction_FunctionShouldPanic(t *testing.T) {
	defer func() {
		e := recover()
		fmt.Println(e)
	}()

	_ = TaskFrom("wrong", 1, 1)

	log.Fatalf("function shoulf have been paniced before")
}

func TestTaskFrom_WhenTaskIsFinished_StatusShouldBeAsExpected(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testData := rand.Intn(20)

	fn := func(i int) int {
		time.Sleep(time.Second)
		return i * 2
	}

	task := TaskFrom("multiplier", fn, testData)
	task.Run()
	task.Wait()

	_, e := task.NotifyWhenDone()

	if e == nil {
		log.Fatalln("setting notifier on finished task should return an error")
	}

	if task.Status() != StatusTaskFinished {
		log.Fatalln("status should be StatusTaskFinished")
	}
}
