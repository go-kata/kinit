package kinitq

import "testing"

func TestGlobal(t *testing.T) {
	if Global() != globalInspector {
		t.Fail()
		return
	}
}
