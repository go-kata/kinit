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

func TestMustInvoke__NilBootstrapper(t *testing.T) {
	err := kerror.Try(func() error {
		kinit.MustInvoke(testExecutor{}, nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}
