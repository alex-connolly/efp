package efp

// These errors are reproduced in the README
// Any changes must be duplicated there.

// Errors during prototype generation.
const (
	errAliasNotVisible = "Alias %s discovered in element %s not found."
	errUnclosedArray   = "Array declaration in field %s is incomplete."
)

// Errors during parsing.
const (
	errDuplicateElement  = "Only %d element(s) with key %s permitted in scope %s."
	errDuplicateField    = "Only %d field(s) with key %s permitted in scope %s."
	errInvalidFieldValue = "Value %s does not match regex %s for field %s."
)
