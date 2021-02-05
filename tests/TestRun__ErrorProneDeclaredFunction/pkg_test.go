package pkg

import (
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestRun__ErrorProneDeclaredFunction(t *testing.T) {
	kinit.MustDeclareErrorProne(func() error {
		return kerror.New(kerror.ECustom, "test error")
	})
	err := kinit.Run()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}
