package efp

type field struct {
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
	parent         *element
	key            string
	parameters     []string
	elementAliases map[string][]*element
	elements       map[string][]*element
	fields         map[string][]*field
	fieldAliases   map[string][]*field
}
