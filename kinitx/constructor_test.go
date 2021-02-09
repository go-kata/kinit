package kinitx

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

type testConstructorT1 struct{}

type testConstructorT2 struct{}

type testConstructorT3 struct {
	obj1 *testConstructorT1
	obj2 *testConstructorT2
}

func TestConstructor__FunctionReturningObject(t *testing.T) {
	ctor := MustNewConstructor(func(
		obj1 *testConstructorT1,
		obj2 *testConstructorT2,
	) *testConstructorT3 {
		return &testConstructorT3{obj1, obj2}
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	obj1 := &testConstructorT1{}
	obj2 := &testConstructorT2{}
	o3, dtor, err := ctor.Create(reflect.ValueOf(obj1), reflect.ValueOf(obj2))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	defer dtor.MustDestroy()
	obj3, ok := o3.Interface().(*testConstructorT3)
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

func TestConstructor__FunctionReturningObjectAndError(t *testing.T) {
	ctor := MustNewConstructor(func(
		obj1 *testConstructorT1,
		obj2 *testConstructorT2,
	) (
		*testConstructorT3,
		error,
	) {
		return &testConstructorT3{obj1, obj2}, nil
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	obj1 := &testConstructorT1{}
	obj2 := &testConstructorT2{}
	o3, dtor, err := ctor.Create(reflect.ValueOf(obj1), reflect.ValueOf(obj2))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	defer dtor.MustDestroy()
	obj3, ok := o3.Interface().(*testConstructorT3)
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

func TestConstructor__FunctionReturningDestructibleObjectAndError(t *testing.T) {
	ctor := MustNewConstructor(func(
		obj1 *testConstructorT1,
		obj2 *testConstructorT2,
	) (
		*testConstructorT3,
		kdone.Destructor,
		error,
	) {
		return &testConstructorT3{obj1, obj2}, kdone.Noop, nil
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	obj1 := &testConstructorT1{}
	obj2 := &testConstructorT2{}
	o3, dtor, err := ctor.Create(reflect.ValueOf(obj1), reflect.ValueOf(obj2))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	defer dtor.MustDestroy()
	obj3, ok := o3.Interface().(*testConstructorT3)
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

func TestNewConstructor__Nil(t *testing.T) {
	_, err := NewConstructor(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewConstructor__NilFunction(t *testing.T) {
	_, err := NewConstructor((func() int)(nil))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewConstructor__WrongType(t *testing.T) {
	_, err := NewConstructor(0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewConstructor__WrongSignature(t *testing.T) {
	_, err := NewConstructor(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestConstructor_Create__WrongNumberOfArguments(t *testing.T) {
	ctor := MustNewConstructor(func(
		obj1 *testConstructorT1,
		obj2 *testConstructorT2,
	) *testConstructorT3 {
		return &testConstructorT3{obj1, obj2}
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	_, _, err := ctor.Create(reflect.ValueOf(&testConstructorT1{}))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestConstructor_Create__WrongArgumentType(t *testing.T) {
	ctor := MustNewConstructor(func(
		obj1 *testConstructorT1,
		obj2 *testConstructorT2,
	) *testConstructorT3 {
		return &testConstructorT3{obj1, obj2}
	})
	t.Logf("%+v %+v", ctor.Type(), ctor.Parameters())
	_, _, err := ctor.Create(reflect.ValueOf(&testConstructorT1{}), reflect.ValueOf(0))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNilConstructor_Type(t *testing.T) {
	if (*Constructor)(nil).Type() != nil {
		t.Fail()
		return
	}
}

func TestNilConstructor_Parameters(t *testing.T) {
	if (*Constructor)(nil).Parameters() != nil {
		t.Fail()
		return
	}
}

func TestNilConstructor_Create(t *testing.T) {
	obj, dtor, err := (*Constructor)(nil).Create()
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
