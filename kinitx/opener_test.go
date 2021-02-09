package kinitx

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

type testOpenerT1 struct{}

type testOpenerT2 struct{}

type testOpenerT3 struct {
	cptr *int
	obj1 *testOpenerT1
	obj2 *testOpenerT2
}

func (t3 *testOpenerT3) Close() error {
	*t3.cptr--
	return nil
}

func TestOpener__FunctionReturningObject(t *testing.T) {
	var c int
	defer func() {
		if c != -1 {
			t.Fail()
			return
		}
	}()
	ctor := MustNewOpener(func(cptr *int, obj1 *testOpenerT1, obj2 *testOpenerT2) *testOpenerT3 {
		return &testOpenerT3{cptr, obj1, obj2}
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	obj1 := &testOpenerT1{}
	obj2 := &testOpenerT2{}
	o3, dtor, err := ctor.Create(reflect.ValueOf(&c), reflect.ValueOf(obj1), reflect.ValueOf(obj2))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	defer dtor.MustDestroy()
	obj3, ok := o3.Interface().(*testOpenerT3)
	if !ok {
		t.Logf("%+v", o3)
		t.Fail()
		return
	}
	if obj3.obj1 != obj1 || obj3.obj2 != obj2 {
		t.Fail()
		return
	}
}

func TestOpener__FunctionReturningObjectAndError(t *testing.T) {
	var c int
	defer func() {
		if c != -1 {
			t.Fail()
			return
		}
	}()
	ctor := MustNewOpener(func(cptr *int, obj1 *testOpenerT1, obj2 *testOpenerT2) (*testOpenerT3, error) {
		return &testOpenerT3{cptr, obj1, obj2}, nil
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	obj1 := &testOpenerT1{}
	obj2 := &testOpenerT2{}
	o3, dtor, err := ctor.Create(reflect.ValueOf(&c), reflect.ValueOf(obj1), reflect.ValueOf(obj2))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	defer dtor.MustDestroy()
	obj3, ok := o3.Interface().(*testOpenerT3)
	if !ok {
		t.Logf("%+v", o3)
		t.Fail()
		return
	}
	if obj3.obj1 != obj1 || obj3.obj2 != obj2 {
		t.Fail()
		return
	}
}

func TestNewOpener__Nil(t *testing.T) {
	_, err := NewOpener(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewOpener__NilFunction(t *testing.T) {
	_, err := NewOpener((func() *testOpenerT3)(nil))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewOpener__WrongType(t *testing.T) {
	_, err := NewOpener(0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewOpener__WrongSignature(t *testing.T) {
	_, err := NewOpener(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewOpener__WrongCloser(t *testing.T) {
	defer func() {
		if v := recover(); v != nil {
			t.Logf("%+v", v)
			t.Fail()
			return
		}
	}()
	_, err := NewOpener(func() struct{} { return struct{}{} })
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestOpener_Create__WrongNumberOfArguments(t *testing.T) {
	ctor := MustNewOpener(func(
		obj1 *testOpenerT1,
		obj2 *testOpenerT2,
	) *testOpenerT3 {
		return &testOpenerT3{new(int), obj1, obj2}
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	_, _, err := ctor.Create(reflect.ValueOf(&testOpenerT1{}))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestOpener_Create__WrongArgumentType(t *testing.T) {
	ctor := MustNewOpener(func(
		obj1 *testOpenerT1,
		obj2 *testOpenerT2,
	) *testOpenerT3 {
		return &testOpenerT3{new(int), obj1, obj2}
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	_, _, err := ctor.Create(reflect.ValueOf(&testOpenerT1{}), reflect.ValueOf(0))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNilOpener_Type(t *testing.T) {
	if (*Opener)(nil).Type() != nil {
		t.Fail()
		return
	}
}

func TestNilOpener_Parameters(t *testing.T) {
	if (*Opener)(nil).Parameters() != nil {
		t.Fail()
		return
	}
}

func TestNilOpener_Create(t *testing.T) {
	obj, dtor, err := (*Opener)(nil).Create()
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
