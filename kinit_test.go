package kinit

import "testing"

func TestGlobal(t *testing.T) {
	if Global() != globalContainer {
		t.Fail()
		return
	}
}
