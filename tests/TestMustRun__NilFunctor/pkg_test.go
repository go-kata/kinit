package pkg

import (
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestMustRun__NilFunctor(t *testing.T) {
	err := kerror.Try(func() error {
		kinit.MustRun(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}
