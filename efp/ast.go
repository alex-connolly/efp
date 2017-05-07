package efp

type field struct {
	alias string
	key   string
	value []*fieldValue
}

type fieldValue struct {
	isArray bool
	regex   string
	values  []*fieldValue
	min     int
	max     int
}

type element struct {
	alias          string
	parent         *element
	key            string
	parameters     []fieldValue
	elementAliases map[string][]*element
	elements       map[string][]*element
	fields         map[string][]*field
	fieldAliases   map[string][]*field
}
