package kinit

import (
	"testing"

	"github.com/go-kata/kerror"
)

func TestApply__NilProcessor(t *testing.T) {
	err := Apply(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestMustApply__NilProcessor(t *testing.T) {
	err := kerror.Try(func() error {
		MustApply(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Apply__NilProcessor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Apply(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Apply__ProcessorWithBrokenType(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Apply(testProcessorWithBrokenType{})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestNilContainer_Apply(t *testing.T) {
	err := (*Container)(nil).Apply(newTestProcessor(processTestCounter))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}
