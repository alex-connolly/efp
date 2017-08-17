package efp

import (
	"fmt"
	"testing"

	"github.com/end-r/goutil"
)

func TestIsPrototypeField(t *testing.T) {
	p := createPrototypeParserString("name : string")
	goutil.Assert(t, isPrototypeField(p), "basic field failed")
	p = createPrototypeParserString("name : string!")
	goutil.Assert(t, isPrototypeField(p), "required field failed")
	p = createPrototypeParserString("<name> : string")
	goutil.Assert(t, isPrototypeField(p), "empty corners field failed")
	p = createPrototypeParserString("<3:name> : string")
	goutil.Assert(t, isPrototypeField(p), "minimum field failed")
	p = createPrototypeParserString("<name:3> : string")
	goutil.Assert(t, isPrototypeField(p), "maximum field failed")
	p = createPrototypeParserString("<3:name:3> : string")
	goutil.Assert(t, isPrototypeField(p), "fixed field failed")
	p = createPrototypeParserString("<name|string> : string")
	goutil.Assert(t, isPrototypeField(p), "disjunction field failed")
	p = createPrototypeParserString("<3:name|string> : string")
	goutil.Assert(t, isPrototypeField(p), "minimum disjunction field failed")
	p = createPrototypeParserString("<name|string:3> : string")
	goutil.Assert(t, isPrototypeField(p), "maximum disjunction field failed")
	p = createPrototypeParserString(`<3:name|string|"a-zA-z{5}":3> : string`)
	goutil.Assert(t, isPrototypeField(p), "fixed disjunction field failed")
	p = createPrototypeParserString(`"[a-z]+" : string`)
	goutil.Assert(t, isPrototypeField(p), "regex key failed")
	p = createPrototypeParserString(`<LIMIT:"[a-z]{3}":LIMIT> : [LIMIT:string:LIMIT]`)
	goutil.Assert(t, isPrototypeField(p), "highly aliased failed")
}

func TestIsPrototypeElement(t *testing.T) {
	p := createPrototypeParserString("name {}")
	goutil.Assert(t, isPrototypeElement(p), "basic element failed")
	p = createPrototypeParserString("<name> {}")
	goutil.Assert(t, isPrototypeElement(p), "corners element failed")
	p = createPrototypeParserString("<3:name> {}")
	goutil.Assert(t, isPrototypeElement(p), "minimum element failed")
	p = createPrototypeParserString("<name:3> {}")
	goutil.Assert(t, isPrototypeElement(p), "maximum element failed")
	p = createPrototypeParserString("<3:name:3> {}")
	goutil.Assert(t, isPrototypeElement(p), "fixed element failed")

	p = createPrototypeParserString("<3:name>(){}")
	goutil.Assert(t, isPrototypeElement(p), "minimum empty parameterised element failed")
	p = createPrototypeParserString("<name:3>(){}")
	goutil.Assert(t, isPrototypeElement(p), "maximum empty parameterised element failed")
	p = createPrototypeParserString("<3:name:3>(){}")
	goutil.Assert(t, isPrototypeElement(p), "fixed empty parameterised element failed")

	p = createPrototypeParserString("<3:name>(int, string){}")
	goutil.Assert(t, isPrototypeElement(p), "minimum parameterised element failed")
	p = createPrototypeParserString("<name:3>(int, string){}")
	goutil.Assert(t, isPrototypeElement(p), "maximum parameterised element failed")
	p = createPrototypeParserString("<3:name:3>(int, string){}")
	goutil.Assert(t, isPrototypeElement(p), "fixed parameterised element failed")
}

func TestIsFieldAlias(t *testing.T) {
	p := createPrototypeParserString("alias x = name : y")
	goutil.Assert(t, isFieldAlias(p), "basic field alias failed")
	p = createPrototypeParserString("alias x = name : string!")
	goutil.Assert(t, isFieldAlias(p), "required field alias failed")
	p = createPrototypeParserString("alias x = <name> : string")
	goutil.Assert(t, isFieldAlias(p), "corners field alias failed")
	p = createPrototypeParserString("alias x = <3:name> : string")
	goutil.Assert(t, isFieldAlias(p), "minimum field alias failed")
	p = createPrototypeParserString("alias x = <name:3> : string")
	goutil.Assert(t, isFieldAlias(p), "maximum field alias failed")
	p = createPrototypeParserString("alias x = <3:name:3> : string")
	goutil.Assert(t, isFieldAlias(p), "fixed field alias failed")
	p = createPrototypeParserString("alias x = <3:name|int> : string")
	goutil.Assert(t, isFieldAlias(p), "minimum disjunction field alias failed")
	p = createPrototypeParserString("alias x = <name|int|string:3> : string")
	goutil.Assert(t, isFieldAlias(p), "maximum disjunction field alias failed")
	p = createPrototypeParserString(`alias x = <3:"A-Za-z":3> : string`)
	goutil.Assert(t, isFieldAlias(p), "fixed field alias with regex failed")
	p = createPrototypeParserString(`alias x = <3:"A-Za-z"|int:3> : string`)
	goutil.Assert(t, isFieldAlias(p), "fixed field disjunction failed")
}

func TestIsElementAlias(t *testing.T) {
	p := createPrototypeParserString("alias x = name {}")
	goutil.Assert(t, isElementAlias(p), "basic element alias failed")
	p = createPrototypeParserString("alias x = <name> {}")
	goutil.Assert(t, isElementAlias(p), "corners element alias failed")
	p = createPrototypeParserString("alias x = <3:name> {}")
	goutil.Assert(t, isElementAlias(p), "minimum element alias failed")
	p = createPrototypeParserString("alias x = <name:3> {}")
	goutil.Assert(t, isElementAlias(p), "maximum element alias failed")
	p = createPrototypeParserString("alias x = <3:name:3> {}")
	goutil.Assert(t, isElementAlias(p), "fixed element alias failed")
	p = createPrototypeParserString("alias x = name(int){}")
	goutil.Assert(t, isElementAlias(p), "basic parameterised element alias failed")
	p = createPrototypeParserString("alias x = <name>(int){}")
	goutil.Assert(t, isElementAlias(p), "corners parameterised element alias failed")
	p = createPrototypeParserString("alias x = <3:name>(int){}")
	goutil.Assert(t, isElementAlias(p), "minimum parameterised element alias failed")
	p = createPrototypeParserString("alias x = <name:3>(int){}")
	goutil.Assert(t, isElementAlias(p), "maximum parameterised element alias failed")
	p = createPrototypeParserString("alias x = <3:name:3>(int){}")
	goutil.Assert(t, isElementAlias(p), "fixed parameterised element alias failed")
	p = createPrototypeParserString("alias x = <name|int|string>(int){}")
	goutil.Assert(t, isElementAlias(p), "disjunction parameterised element alias failed")
	p = createPrototypeParserString("alias x = <3:name|int|string>(int){}")
	goutil.Assert(t, isElementAlias(p), "minimum disjunction parameterised element alias failed")
	p = createPrototypeParserString("alias x = <name|int|string:3>(int){}")
	goutil.Assert(t, isElementAlias(p), "maximum disjunction parameterised element alias failed")
	p = createPrototypeParserString("alias x = <3:name|int|string:3>(int){}")
	goutil.Assert(t, isElementAlias(p), "disjunction parameterised element alias failed")
	p = createPrototypeParserString(`alias x = <3:name|"[a-zA-Z]+"|string:3>(int){}`)
	goutil.Assert(t, isElementAlias(p), "regex disjunction parameterised element alias failed")
}

func TestIsField(t *testing.T) {
	p := createPrototypeParserString("name = 6")
	goutil.Assert(t, isField(p), "int field failed")
	p = createPrototypeParserString(`name = "www"`)
	goutil.Assert(t, isField(p), "string field failed")
	p = createPrototypeParserString("name = hi")
	goutil.Assert(t, isField(p), "id field failed")
	p = createPrototypeParserString(`name = [hi, me, c]`)
	goutil.Assert(t, isField(p), "array field failed")
	p = createPrototypeParserString(`name = ["hi", "me", "c"]`)
	goutil.Assert(t, isField(p), "string array field failed")
	p = createPrototypeParserString(`name = [["hi", "me"], [6, 5, 7]]`)
	goutil.Assert(t, isField(p), "2D array field failed")
}

func TestIsElement(t *testing.T) {
	p := createPrototypeParserString("name {}")
	goutil.Assert(t, isElement(p), "basic element failed")
	p = createPrototypeParserString(`"name" {}`)
	goutil.Assert(t, isElement(p), "basic regex element failed")
	p = createPrototypeParserString("name(int){}")
	goutil.Assert(t, isElement(p), "parameterised element failed")
	p = createPrototypeParserString("name(){}")
	goutil.Assert(t, isElement(p), "empty parameterised element failed")
	p = createPrototypeParserString(`"name"(){}`)
	goutil.Assert(t, isElement(p), "empty regex parameterised element failed")
	p = createPrototypeParserString("name(int, int, string){}")
	goutil.Assert(t, isElement(p), "multi parameterised element failed")
}

func TestIsOperator(t *testing.T) {
	// doesn't really need a test
	ops := []byte{',', '|', '<', '>', '{', '}', '!', '[', ']', '(', ')'}
	for _, op := range ops {
		goutil.Assert(t, is(op)(op), "operator failed")
	}
}

func TestIsDistant(t *testing.T) {

	p := createPrototypeParserString("alias x = <3:name|int|string:3>(int){}")
	goutil.Assert(t, realDistance(p, tknOpenCorner, 1) == 3,
		fmt.Sprintf("wrong corner distance: %d", realDistance(p, tknOpenCorner, 1)))
	goutil.Assert(t, realDistance(p, tknOpenBracket, 1) == 6,
		fmt.Sprintf("wrong bracket distance: %d", realDistance(p, tknOpenBracket, 1)))

	p = createPrototypeParserString("alias x = <3: name | int | string :3>(int){}")
	goutil.Assert(t, realDistance(p, tknOpenCorner, 1) == 3,
		fmt.Sprintf("wrong corner distance: %d", realDistance(p, tknOpenCorner, 1)))
	goutil.Assert(t, realDistance(p, tknOpenBracket, 1) == 6,
		fmt.Sprintf("wrong bracket distance: %d", realDistance(p, tknOpenBracket, 1)))

	p = createPrototypeParserString(`alias x = <3: name | "[A-Z]+"| string :3>(int){}`)
	goutil.Assert(t, realDistance(p, tknOpenCorner, 1) == 3,
		fmt.Sprintf("wrong corner distance: %d", realDistance(p, tknOpenCorner, 1)))
	goutil.Assert(t, realDistance(p, tknOpenBracket, 1) == 6,
		fmt.Sprintf("wrong bracket distance: %d", realDistance(p, tknOpenBracket, 1)))

	p = createPrototypeParserString("<3:name> : string")
	goutil.Assert(t, realDistance(p, tknOpenCorner, 1) == 0,
		fmt.Sprintf("wrong corner distance: %d", realDistance(p, tknOpenCorner, 1)))
	goutil.Assert(t, realDistance(p, tknColon, 1) == 3,
		fmt.Sprintf("wrong colon distance: %d", realDistance(p, tknColon, 1)))

	//empty bytes
	p = createPrototypeParserString("")
	goutil.Assert(t, realDistance(p, tknValue, 1) == -1, "failed empty")

	p = createPrototypeParserString("ALIAS1 ALIAS2")
	goutil.Assert(t, realDistance(p, tknValue, 1) == 0, "failed alias")

}

func TestIsTextAlias(t *testing.T) {
	p := createPrototypeParserString(`alias x = hi`)
	goutil.Assert(t, isValueAlias(p), "id text alias failed")
	p = createPrototypeParserString(`alias x = "5"`)
	goutil.Assert(t, isValueAlias(p), "string text alias failed")
}

func TestIsDiscoveredAlias(t *testing.T) {
	p := createPrototypeParserString("hello")
	goutil.Assert(t, isDiscoveredAlias(p), "discovered alias failed")
}
