package efp

// These errors are reproduced in the README
// Any changes must be duplicated there.

// Errors during prototype generation.
const (
	errDuplicateAlias        = "Alias %s already declared in scope %s."
	errAliasNotVisible       = "Alias %s discovered in element %s not found."
	errUnclosedArray         = "Array declaration in field %s is incomplete."
	errInvalidRegex          = "Invalid regex string %s in element %s."
	errInvalidLimitAlias     = "Invalid limit alias."
	errUnrecognisedConstruct = "Unrecognised construct %s."
	errUnrecognisedToken     = "Unrecognised token %s."
)

// Errors during parsing.
const (
	errDuplicateElement    = "Only %d element(s) with key %s permitted in scope %s (found %d)."
	errDuplicateField      = "Only %d field(s) with key %s permitted in scope %s (found %d)."
	errInvalidFieldValue   = "Value %s does not match regex %s for field %s."
	errInvalidToken        = "Invalid token %s in %s."
	errInsufficientFields  = "%d field(s) with key %s discovered in element %s (%d required)."
	errArrayMinimum        = "Field array %s in scope %s must have at least %d values (found %d)."
	errArrayMaximum        = "Field array %s in scope %s must not have more than %d values (found %d)."
	errUnmatchedFieldValue = "Invalid value for field key '%s'."
	errRequiredArray       = "Array type required."
)

var standards = map[string]TextAlias{
	"string": TextAlias{`^(.*)$`, false},
	"int":    TextAlias{`^([-]?[1-9]\d*|0)$`, false},
	"float":  TextAlias{"^(([-]?[1-9][0-9]*)|([-]?[0-9]+.[0-9]+)|0)$", false},
	"bool":   TextAlias{"^(true|false)$", false},
	"uint":   TextAlias{`^([1-9]\d*|0)$`, false},
}
