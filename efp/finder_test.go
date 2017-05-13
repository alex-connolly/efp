package efp

import (
	"fmt"
	"testing"
)

func TestIsPrototypeField(t *testing.T) {
	p := &parser{lexer: lex([]byte("name : string"))}
	assert(t, isPrototypeField(p), "basic field failed")
	p = &parser{lexer: lex([]byte("name : string!"))}
	assert(t, isPrototypeField(p), "required field failed")
	p = &parser{lexer: lex([]byte("<name> : string"))}
	assert(t, isPrototypeField(p), "empty corners field failed")
	p = &parser{lexer: lex([]byte("<3:name> : string"))}
	assert(t, isPrototypeField(p), "minimum field failed")
	p = &parser{lexer: lex([]byte("<name:3> : string"))}
	assert(t, isPrototypeField(p), "maximum field failed")
	p = &parser{lexer: lex([]byte("<3:name:3> : string"))}
	assert(t, isPrototypeField(p), "fixed field failed")
	p = &parser{lexer: lex([]byte("<name|string> : string"))}
	assert(t, isPrototypeField(p), "disjunction field failed")
	p = &parser{lexer: lex([]byte("<3:name|string> : string"))}
	assert(t, isPrototypeField(p), "minimum disjunction field failed")
	p = &parser{lexer: lex([]byte("<name|string:3> : string"))}
	assert(t, isPrototypeField(p), "maximum disjunction field failed")
	p = &parser{lexer: lex([]byte(`<3:name|string|"a-zA-z{5}":3> : string`))}
	assert(t, isPrototypeField(p), "fixed disjunction field failed")
}

func TestIsPrototypeElement(t *testing.T) {
	p := &parser{lexer: lex([]byte("name {}"))}
	assert(t, isPrototypeElement(p), "basic element failed")
	p = &parser{lexer: lex([]byte("<name> {}"))}
	assert(t, isPrototypeElement(p), "corners element failed")
	p = &parser{lexer: lex([]byte("<3:name> {}"))}
	assert(t, isPrototypeElement(p), "minimum element failed")
	p = &parser{lexer: lex([]byte("<name:3> {}"))}
	assert(t, isPrototypeElement(p), "maximum element failed")
	p = &parser{lexer: lex([]byte("<3:name:3> {}"))}
	assert(t, isPrototypeElement(p), "fixed element failed")

	p = &parser{lexer: lex([]byte("<3:name>(){}"))}
	assert(t, isPrototypeElement(p), "minimum empty parameterised element failed")
	p = &parser{lexer: lex([]byte("<name:3>(){}"))}
	assert(t, isPrototypeElement(p), "maximum empty parameterised element failed")
	p = &parser{lexer: lex([]byte("<3:name:3>(){}"))}
	assert(t, isPrototypeElement(p), "fixed emptyparameterised element failed")

	p = &parser{lexer: lex([]byte("<3:name>(int, string){}"))}
	assert(t, isPrototypeElement(p), "minimum parameterised element failed")
	p = &parser{lexer: lex([]byte("<name:3>(int, string){}"))}
	assert(t, isPrototypeElement(p), "maximum parameterised element failed")
	p = &parser{lexer: lex([]byte("<3:name:3>(int, string){}"))}
	assert(t, isPrototypeElement(p), "fixed parameterised element failed")
}

func TestIsFieldAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("alias x = name : string"))}
	assert(t, isFieldAlias(p), "basic field alias failed")
	p = &parser{lexer: lex([]byte("alias x = name : string!"))}
	assert(t, isFieldAlias(p), "required field alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name> : string"))}
	assert(t, isFieldAlias(p), "corners field alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name> : string"))}
	assert(t, isFieldAlias(p), "minimum field alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name:3> : string"))}
	assert(t, isFieldAlias(p), "maximum field alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name:3> : string"))}
	assert(t, isFieldAlias(p), "fixed field alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name|int> : string"))}
	assert(t, isFieldAlias(p), "minimum disjunction field alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name|int|string:3> : string"))}
	assert(t, isFieldAlias(p), "maximum disjunction field alias failed")
	p = &parser{lexer: lex([]byte(`alias x = <3:"A-Za-z":3> : string`))}
	assert(t, isFieldAlias(p), "fixed field alias with regex failed")
	p = &parser{lexer: lex([]byte(`alias x = <3:"A-Za-z"|int:3> : string`))}
	assert(t, isFieldAlias(p), "fixed field disjunction failed")
}

func TestIsElementAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("alias x = name {}"))}
	assert(t, isElementAlias(p), "basic element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name> {}"))}
	assert(t, isElementAlias(p), "corners element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name> {}"))}
	assert(t, isElementAlias(p), "minimum element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name:3> {}"))}
	assert(t, isElementAlias(p), "maximum element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name:3> {}"))}
	assert(t, isElementAlias(p), "fixed element alias failed")
	p = &parser{lexer: lex([]byte("alias x = name(int){}"))}
	assert(t, isElementAlias(p), "basic parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name>(int){}"))}
	assert(t, isElementAlias(p), "corners parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name>(int){}"))}
	assert(t, isElementAlias(p), "minimum parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name:3>(int){}"))}
	assert(t, isElementAlias(p), "maximum parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name:3>(int){}"))}
	assert(t, isElementAlias(p), "fixed parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name|int|string>(int){}"))}
	assert(t, isElementAlias(p), "disjunction parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name|int|string>(int){}"))}
	assert(t, isElementAlias(p), "minimum disjunction parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <name|int|string:3>(int){}"))}
	assert(t, isElementAlias(p), "maximum disjunction parameterised element alias failed")
	p = &parser{lexer: lex([]byte("alias x = <3:name|int|string:3>(int){}"))}
	assert(t, isElementAlias(p), "disjunction parameterised element alias failed")
	p = &parser{lexer: lex([]byte(`alias x = <3:name|"a-zA-Z"|string:3>(int){}`))}
	assert(t, isElementAlias(p), "regex disjunction parameterised element alias failed")
}

func TestIsField(t *testing.T) {

	p := &parser{lexer: lex([]byte("name = 6"))}
	assert(t, isField(p), "int field failed")
	p = &parser{lexer: lex([]byte(`name = "www"`))}
	assert(t, isField(p), "string field failed")
	p = &parser{lexer: lex([]byte("name = hi"))}
	assert(t, isField(p), "id field failed")
	p = &parser{lexer: lex([]byte(`name = [hi, me, c]`))}
	assert(t, isField(p), "array field failed")
	p = &parser{lexer: lex([]byte(`name = ["hi", "me", "c"]`))}
	assert(t, isField(p), "string array field failed")
	p = &parser{lexer: lex([]byte(`name = [["hi", "me"], [6, 5, 7]]`))}
	assert(t, isField(p), "2D array field failed")
}

func TestIsElement(t *testing.T) {
	p := &parser{lexer: lex([]byte("name {}"))}
	assert(t, isElement(p), "basic element failed")
	p = &parser{lexer: lex([]byte("name(int){}"))}
	assert(t, isElement(p), "parameterised element failed")
	p = &parser{lexer: lex([]byte("name(){}"))}
	assert(t, isElement(p), "empty parameterised element failed")
	p = &parser{lexer: lex([]byte("name(int, int, string){}"))}
	assert(t, isElement(p), "multi parameterised element failed")
}

func TestIsOperator(t *testing.T) {
	// doesn't really need a test
	ops := []byte{',', '|', '<', '>', '{', '}', '!', '[', ']', '(', ')'}
	for _, op := range ops {
		assert(t, is(op)(op), "operator failed")
	}
}

func TestIsDistant(t *testing.T) {

	p := &parser{lexer: lex([]byte("alias x = <3:name|int|string:3>(int){}"))}
	assert(t, realDistance(p, tknOpenCorner, 1) == 3,
		fmt.Sprintf("wrong corner distance: %d", realDistance(p, tknOpenCorner, 1)))
	assert(t, realDistance(p, tknOpenBracket, 1) == 10,
		fmt.Sprintf("wrong bracket distance: %d", realDistance(p, tknOpenBracket, 1)))

	p = &parser{lexer: lex([]byte("alias x = <3: name | int | string :3>(int){}"))}
	assert(t, realDistance(p, tknOpenCorner, 1) == 3,
		fmt.Sprintf("wrong corner distance: %d", realDistance(p, tknOpenCorner, 1)))
	assert(t, realDistance(p, tknOpenBracket, 1) == 10,
		fmt.Sprintf("wrong bracket distance: %d", realDistance(p, tknOpenBracket, 1)))

	// empty bytes
	p = &parser{lexer: lex([]byte(""))}
	assert(t, realDistance(p, tknValue, 1) == -1, "failed empty")

	p = &parser{lexer: lex([]byte("ALIAS1 ALIAS2"))}
	assert(t, realDistance(p, tknValue, 1) == 0, "failed alias")

}

func TestIsTextAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("alias x = 5"))}
	assert(t, isTextAlias(p), "int text alias failed")
	p = &parser{lexer: lex([]byte(`alias x = hi`))}
	assert(t, isTextAlias(p), "id text alias failed")
	p = &parser{lexer: lex([]byte(`alias x = "5"`))}
	assert(t, isTextAlias(p), "string text alias failed")
}

func TestIsAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("alias x ="))}
	assert(t, isAlias(p), "not an alias")
}

func TestIsDiscoveredAlias(t *testing.T) {
	p := &parser{lexer: lex([]byte("hello"))}
	assert(t, isDiscoveredAlias(p), "discovered alias failed")
}
