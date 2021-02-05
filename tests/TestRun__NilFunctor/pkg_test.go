package pkg

import (
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestRun__NilFunctor(t *testing.T) {
	err := kinit.Run(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}
