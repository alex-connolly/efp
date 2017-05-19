package efp

import "regexp"

type field struct {
	alias string
	key   *key
	value *value
}

type key struct {
	key   string
	regex *regexp.Regexp
	min   int
	max   int
}

type value struct {
	parent   *value
	value    string
	children []*value
}

type typeDeclaration struct {
	isArray bool
	types   []*typeDeclaration
	value   string
	min     int
	max     int
}
type protoField struct {
	key   *key
	types []*typeDeclaration
}

type protoElement struct {
	alias          string
	key            *key
	parent         *protoElement
	parameters     []*typeDeclaration
	fields         map[string]*protoField
	fieldAliases   map[string]*protoField
	elements       map[string]*protoElement
	elementAliases map[string]*protoElement
	aliases        []string
	textAliases    map[string]string
}

type element struct {
	alias      string
	key        *key
	parent     *element
	parameters []*value
	elements   map[string][]*element
	fields     map[string][]*field
}

var standards = map[string]string{
	"string": `^(.*)$`,
	"int":    `^([-]?[1-9]\d*|0)$`,
	"float":  "^([-]?([0-9]*[.])?[0-9]+|)$",
	"bool":   "^(true|false)$",
	"uint":   `^([1-9]\d*|0)$`,
}

func (p *protoElement) addStandardAliases() {
	p.textAliases = standards
}
