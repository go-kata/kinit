package kinitq

import (
	"reflect"
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestInspector__OK(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func() int16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func() int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int16) int64 { return 0 }))
	ctr.MustAttach(newTestProcessor(func(int64, int32) {}))
	ctr.MustProvide(newTestConstructor(func(int32) uint64 { return 0 }))
	inspector := NewInspector()
	inspector.MustRequire(reflect.TypeOf(int64(0)))
	inspector.MustRequire(reflect.TypeOf(uint64(0)))
	if err := inspector.Inspect(ctr, nil); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestInspector__OKWhenIgnore(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(string) int16 { return 0 })) // unsatisfied dependency: string
	inspector := NewInspector()
	inspector.MustIgnore(reflect.TypeOf(""))
	if err := inspector.Inspect(ctr, nil); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestInspector__OKWhenInspectOnlyRequired(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(string) int16 { return 0 })) // unsatisfied dependency: string
	ctr.MustProvide(newTestConstructor(func() int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func() int64 { return 0 }))
	ctr.MustAttach(newTestProcessor(func(int64, int32) {}))
	ctr.MustProvide(newTestConstructor(func(int32) uint64 { return 0 }))
	inspector := NewInspector()
	inspector.MustRequire(reflect.TypeOf(int64(0)))
	inspector.MustRequire(reflect.TypeOf(uint64(0)))
	if err := inspector.Inspect(ctr, &Options{
		InspectOnlyRequired: true,
	}); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestInspector__OKWhenAllowIrrelevantProcessors(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustAttach(newTestProcessor(func(int64) {}))
	if err := NewInspector().Inspect(ctr, &Options{
		AllowIrrelevantProcessors: true,
	}); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestInspector__CyclicRequiredDependency(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(int64) int16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int16) int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int32) int64 { return 0 }))
	inspector := NewInspector()
	inspector.MustRequire(reflect.TypeOf(int64(0)))
	err := inspector.Inspect(ctr, &Options{
		InspectOnlyRequired: true,
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestInspector__CyclicConstructorDependency(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(int64) int16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int16) int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int32) int64 { return 0 }))
	err := NewInspector().Inspect(ctr, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestInspector__CyclicProcessorDependency(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(int64) int16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int16) int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func() int64 { return 0 }))
	ctr.MustAttach(newTestProcessor(func(int64, int32) {}))
	err := NewInspector().Inspect(ctr, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestInspector__TwoCyclicDependencies(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(int64) int16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int16) int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int64) uint16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(uint16) uint32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int32, uint32) int64 { return 0 }))
	inspector := NewInspector()
	err := inspector.Inspect(ctr, nil)
	t.Logf("%+v", err)
	errs, ok := err.(kerror.MultiError)
	if !ok {
		t.Fail()
		return
	}
	if len(errs) != 2 || kerror.ClassOf(errs[0]) != kerror.EAmbiguous || kerror.ClassOf(errs[1]) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestInspector__UnsatisfiedRequiredDependency(t *testing.T) {
	ctr := kinit.NewContainer()
	inspector := NewInspector()
	inspector.MustRequire(reflect.TypeOf(int64(0)))
	err := inspector.Inspect(ctr, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENotFound {
		t.Fail()
		return
	}
}

func TestInspector__UnsatisfiedConstructorDependency(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(int32) int64 { return 0 }))
	err := NewInspector().Inspect(ctr, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENotFound {
		t.Fail()
		return
	}
}

func TestInspector__UnsatisfiedProcessorDependency(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func() int64 { return 0 }))
	ctr.MustAttach(newTestProcessor(func(int64, int32) {}))
	err := NewInspector().Inspect(ctr, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENotFound {
		t.Fail()
		return
	}
}

func TestInspector__IrrelevantProcessors(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustAttach(newTestProcessor(func(int64) {}))
	err := NewInspector().Inspect(ctr, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestInspector__MultipleBreaks(t *testing.T) {
	ctr := kinit.NewContainer()
	ctr.MustProvide(newTestConstructor(func(int64) int16 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int16) int32 { return 0 }))
	ctr.MustProvide(newTestConstructor(func(int32) int64 { return 0 }))
	ctr.MustAttach(newTestProcessor(func(uintptr, string) {}))
	inspector := NewInspector()
	inspector.MustRequire(reflect.TypeOf(int64(0)))
	inspector.MustRequire(reflect.TypeOf((*kerror.Collector)(nil)))
	err := inspector.Inspect(ctr, nil)
	t.Logf("%+v", err)
	if !kerror.Is(err, kerror.EAmbiguous) || !kerror.Is(err, kerror.ENotFound) || !kerror.Is(err, kerror.EInvalid) {
		t.Fail()
		return
	}
}

func TestInspector_Require__NilType(t *testing.T) {
	inspector := NewInspector()
	err := inspector.Require(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestInspector_Ignore__NilType(t *testing.T) {
	inspector := NewInspector()
	err := inspector.Ignore(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestNilInspector_Require(t *testing.T) {
	err := (*Inspector)(nil).Require(reflect.TypeOf(0))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilInspector_Ignore(t *testing.T) {
	err := (*Inspector)(nil).Ignore(reflect.TypeOf(0))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilInspector_Inspect(t *testing.T) {
	if err := (*Inspector)(nil).Inspect(kinit.NewContainer(), nil); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}
