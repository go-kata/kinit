package pkg

import (
	"reflect"
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

type testExecutor struct{}

func (testExecutor) Parameters() []reflect.Type {
	return nil
}

func (testExecutor) Execute(a ...reflect.Value) (kinit.Executor, error) {
	return nil, nil
}

func TestInvoke__NilBootstrapper(t *testing.T) {
	err := kinit.Invoke(testExecutor{}, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}
