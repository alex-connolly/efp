package efp

import "testing"

func TestIsOperator(t *testing.T) {
	if !is(',')(',') {
		t.Fail()
	}
}
