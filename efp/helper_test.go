package efp

import "testing"

func assert(t *testing.T, condition bool, err string) {
	if condition {
		t.Log(err)
		t.Fail()
	}
}

func assertNow(t *testing.T, condition bool, err string) {
	if condition {
		t.Log(err)
		t.FailNow()
	}
}
