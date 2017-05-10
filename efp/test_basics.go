package efp

import "testing"

func failIf(t *testing.T, condition bool) {
	if condition {
		t.Fail()
	}
}
