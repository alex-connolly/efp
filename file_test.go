package efp

import (
	"fmt"
	"testing"

	"github.com/end-r/goutil"
)

func testFile(name string) string {
	return fmt.Sprintf("test_files/%s", name)
}

func TestUnknownFile(t *testing.T) {
	_, errs := PrototypeFile("not_found.efp")
	goutil.Assert(t, errs != nil, "errs should not be nil")
}

func TestEmptyFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("empty.efp"))
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
}

func TestLargeFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("large.efp"))
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
}

func TestLargeFiles(t *testing.T) {
	p, errs := PrototypeFile(testFile("large.efp"))
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
	//e, errs := p.ValidateFiles("")
}

func TestVMGenFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("vm.efp"))
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
	e, errs := p.ValidateString(`
		name = "Example"
		author = "[7][7][7]"

		instruction("ADD", "01"){
		    description = "Finds the sum of two numbers."
		    fuel = 100
		}

		instruction("PUSH", "02"){
		    description = "Pushes a number onto the stack."
		    fuel = 30
		}

		instruction("TEST", "03"){
		    description = "Test instruction."
		    fuel = 30
		}
	`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
	goutil.Assert(t, len(e.Elements("instruction")) == 3, "wrong instruction length")
	goutil.Assert(t, e.FirstField("name").Value(0) == "Example",
		fmt.Sprintf("wrong param value %s, expected %s\n",
			e.FirstField("name").Value(0), "Example"))
	goutil.Assert(t, e.FirstElement("instruction").Parameter(0).Value() == "ADD",
		fmt.Sprintf("wrong param value %s, expected %s\n",
			e.FirstElement("instruction").Parameter(0).Value(), "ADD"))
}

func TestComplexVMGenFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("complex_vm.efp"))
	goutil.Assert(t, errs == nil, "prototype errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
	goutil.Assert(t, len(p.elements) == 2, "should be two top level prototype elements")
	goutil.Assert(t, len(p.fields) == 2, "should be two top level prototype elements")
	goutil.Assert(t, p.Element("category") != nil, "should be a top level category element")
	goutil.Assert(t, p.Element("instruction") != nil, "should be a top level category element")
	goutil.Assert(t, p.Element("category").Element("instruction") != nil, "should be a second level category element")
	e, errs := p.ValidateString(`
		name = "Example"
		author = "[7][7][7]"

		instruction("ADD", "01"){
		    description = "Finds the sum of two numbers."
			validate = 2
			fuel = 100
		}

		category("Pushers"){

			description = "Things which push in general."

		    instruction("PUSH", "02"){
		        description = "Pushes a number onto the stack."
				validate = 0
		        fuel = 10
		    }

		    instruction("TEST", "03"){
		        description = "Test instruction."
				validate = 0
		        fuel = 30
		    }

		}

	`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")
}

func TestMinMaxFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("min_max.efp"))
	goutil.Assert(t, errs == nil, "prototype errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
	goutil.Assert(t, len(p.elements) == 1, "should be one top level prototype element")
	goutil.Assert(t, len(p.fields) == 0, "should be zero top level prototype elements")
	goutil.Assert(t, p.Element("dog") != nil, "should be a top level dog element")
	goutil.Assert(t, p.Element("dog").Field("name") != nil, "should be a second level field element")
	e, errs := p.ValidateString(`
		dog("Alex"){
			name = "AC"
			name = "Centurion"
		}
	`)
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, e != nil, "e should not be nil")

}

func TestMinMaxFileErrors(t *testing.T) {
	p, errs := PrototypeFile(testFile("min_max.efp"))
	goutil.Assert(t, errs == nil, "prototype errs should be nil")
	goutil.AssertNow(t, p != nil, "p should not be nil")
	_, errs = p.ValidateString(`
		dog("Alex"){
			name = "AC"
		}
	`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
	_, errs = p.ValidateString(`
		dog("Alex"){
			name = "AC"
			name = "AC"
			name = "AC"
			name = "AC"
			name = "AC"
		}
	`)
	goutil.Assert(t, errs != nil, "errs should not be nil")
}
