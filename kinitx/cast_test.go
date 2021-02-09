package kinitx

import (
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

type testCloser struct{}

func (*testCloser) Close() error {
	return nil
}

func TestCastToConstructor_Interface(t *testing.T) {
	orig := MustNewConstructor(func() int { return 0 })
	ctor, err := castToConstructor(orig)
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if ctor != orig {
		t.Fail()
		return
	}
}

func TestCastToConstructor__Constructor(t *testing.T) {
	ctor, err := castToConstructor(func() int { return 0 })
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := ctor.(*Constructor); !ok {
		t.Fail()
		return
	}
}

func TestCastToConstructor__Opener(t *testing.T) {
	ctor, err := castToConstructor(func() *testCloser { return &testCloser{} })
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := ctor.(*Opener); !ok {
		t.Fail()
		return
	}
}

func TestCastToConstructor__ErrorProneOpener(t *testing.T) {
	ctor, err := castToConstructor(func() (*testCloser, error) { return &testCloser{}, nil })
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := ctor.(*Opener); !ok {
		t.Fail()
		return
	}
}

func TestCastToConstructor__ConstructorOfCloser(t *testing.T) {
	ctor, err := castToConstructor(func() (*testCloser, kdone.Destructor, error) {
		return &testCloser{}, kdone.Noop, nil
	})
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := ctor.(*Constructor); !ok {
		t.Fail()
		return
	}
}

func TestCastToConstructor__Struct(t *testing.T) {
	ctor, err := castToConstructor(testCloser{})
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := ctor.(*Initializer); !ok {
		t.Fail()
		return
	}
}

func TestCastToConstructor__StructPointer(t *testing.T) {
	ctor, err := castToConstructor((*testCloser)(nil))
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := ctor.(*Initializer); !ok {
		t.Fail()
		return
	}
}

func TestCastToConstructor__Nil(t *testing.T) {
	_, err := castToConstructor(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestCastToConstructor__WrongType(t *testing.T) {
	_, err := castToConstructor(0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestCastToConstructor__WrongFunc(t *testing.T) {
	_, err := castToConstructor(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestCastToConstructor__WrongOpener(t *testing.T) {
	_, err := castToConstructor(func() (*testCloser, *testCloser) {
		return &testCloser{}, &testCloser{}
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestCastToConstructor__WrongPointer(t *testing.T) {
	_, err := castToConstructor((*int)(nil))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestCastToProcessor__Interface(t *testing.T) {
	orig := MustNewProcessor(func(int) {})
	proc, err := castToProcessor(orig)
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if proc != orig {
		t.Fail()
		return
	}
}

func TestCastToProcessor__Processor(t *testing.T) {
	proc, err := castToProcessor(func(int) {})
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := proc.(*Processor); !ok {
		t.Fail()
		return
	}
}

func TestCastToProcessor__Nil(t *testing.T) {
	_, err := castToProcessor(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestCastToFunctor__Interface(t *testing.T) {
	orig := MustNewFunctor(func() {})
	fun, err := castToFunctor(orig)
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if fun != orig {
		t.Fail()
		return
	}
}

func TestCastToFunctor__Functor(t *testing.T) {
	fun, err := castToFunctor(func() {})
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := fun.(*Functor); !ok {
		t.Fail()
		return
	}
}

func TestCastToFunctor__Injector(t *testing.T) {
	fun, err := castToFunctor(1)
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if _, ok := fun.(*Injector); !ok {
		t.Fail()
		return
	}
}

func TestCastToFunctor__Nil(t *testing.T) {
	_, err := castToFunctor(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}
