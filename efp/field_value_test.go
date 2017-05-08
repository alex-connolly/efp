package efp

import "testing"

func TestComplexValues(t *testing.T) {
	//Prototype("../samples/complex.efp")
	//	assert(t, p.prototype.fields["name"])
}

func assert(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fail()
	}
}
