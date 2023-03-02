package future_test

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/powerpuffpenguin/async/future"
)

func TestValue(t *testing.T) {
	for i := 0; i < 10; i++ {
		v, e := future.Value(i).Get()
		if e != nil {
			t.Fatal(`value err:`, e)
		}
		if v != i {
			t.Fatalf(`expect=%v actual=%v`, i, v)
		}
	}
}
func TestErr(t *testing.T) {
	var e0 = errors.New(`test err value`)
	v, e1 := future.Error[any](e0).Get()
	if v != nil {
		t.Fatal(`unexpected value:`, v)
	}
	if e1 != e0 {
		t.Fatalf(`expect=%v actual=%v`, e0, e1)
	}
}
func TestCompleter(t *testing.T) {
	var w sync.WaitGroup
	c := future.New[int]()
	w.Add(1)
	var e0 error
	go func() {
		defer w.Done()
		if !c.Complete(1) {
			e0 = errors.New(`Complete(1) false`)
			return
		}
		if c.Complete(2) {
			e0 = errors.New(`Complete(2) true`)
		}
	}()
	v, e1 := c.Future().Get()
	if e1 != nil {
		t.Fatal(`future get err:`, e1)
	} else if v != 1 {
		t.Fatalf(`future value: expect=1 actual=%v`, v)
	}
	w.Wait()
	if e0 != nil {
		t.Fatal(`e0:`, e0)
	}
}
func BenchmarkRead(b *testing.B) {
	var wait sync.WaitGroup
	wait.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wait.Done()
			os.ReadFile("future.go")
		}()
	}
	wait.Wait()
}
func BenchmarkReadFuture(b *testing.B) {
	once := future.NewOnce(func() ([]byte, error) {
		return os.ReadFile("future.go")
	})

	for i := 0; i < b.N; i++ {
		once.Do()
	}
	once.Do().Get()
}

func BenchmarkExecutor(b *testing.B) {
	once := future.NewOnceExecute(func() ([]byte, error) {
		return os.ReadFile("future.go")
	})
	executor := future.NewExecutor()
	for i := 0; i < b.N; i++ {
		executor.Do(once)
	}
	executor.Do(once).Get()
}
