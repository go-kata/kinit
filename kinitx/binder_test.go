package kinitx

import (
	"io"
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

func TestBinder(t *testing.T) {
	ctor := MustNewBinder((*io.Closer)(nil), (kdone.CloserFunc)(nil))
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	var c int
	obj, dtor, err := ctor.Create(reflect.ValueOf(kdone.CloserFunc(func() error {
		c++
		return nil
	})))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if obj.Type() != closerType {
		t.Fail()
		return
	}
	closer, ok := obj.Interface().(io.Closer)
	if !ok {
		t.Fail()
		return
	}
	if err := closer.Close(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
	f, ok := dtor.(kdone.DestructorFunc)
	if !ok {
		t.Fail()
		return
	}
	if err := f.Destroy(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestNewBinder__NilInterfacePointer(t *testing.T) {
	_, err := NewBinder(nil, 0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewBinder__NilObject(t *testing.T) {
	_, err := NewBinder((*io.Closer)(nil), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewBinder__WrongInterfacePointer1(t *testing.T) {
	_, err := NewBinder(0, 0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewBinder__WrongInterfacePointer2(t *testing.T) {
	_, err := NewBinder((*int)(nil), 0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewBinder__IncompatibleObject(t *testing.T) {
	_, err := NewBinder((*io.Closer)(nil), 0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestBinder_Create__WrongArgumentNumber(t *testing.T) {
	ctor := MustNewBinder((*io.Closer)(nil), (kdone.CloserFunc)(nil))
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	_, _, err := ctor.Create()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestBinder_Create__WrongArgumentType(t *testing.T) {
	ctor := MustNewBinder((*io.Closer)(nil), (kdone.CloserFunc)(nil))
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	_, _, err := ctor.Create(reflect.ValueOf(0))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNilBinder_Type(t *testing.T) {
	if (*Binder)(nil).Type() != nil {
		t.Fail()
		return
	}
}

func TestNilBinder_Parameters(t *testing.T) {
	if (*Binder)(nil).Parameters() != nil {
		t.Fail()
		return
	}
}

func TestNilBinder_Create(t *testing.T) {
	obj, dtor, err := (*Binder)(nil).Create()
	if obj != reflect.ValueOf(nil) {
		t.Fail()
		return
	}
	f, ok := dtor.(kdone.DestructorFunc)
	if !ok {
		t.Fail()
		return
	}
	if err := f.Destroy(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}
