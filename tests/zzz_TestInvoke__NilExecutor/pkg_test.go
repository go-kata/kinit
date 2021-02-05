package pkg

import (
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestInvoke__NilExecutor(t *testing.T) {
	err := kinit.Invoke(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}
