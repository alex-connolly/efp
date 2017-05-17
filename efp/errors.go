package efp

// These errors are reproduced in the README
// Any changes must be duplicated there.

// Errors during prototype generation.
const (
	errDuplicateAlias  = "Alias %s already declared in scope %s."
	errAliasNotVisible = "Alias %s discovered in element %s not found."
	errUnclosedArray   = "Array declaration in field %s is incomplete."
	errInvalidRegex    = "Invalid regex string %s in element %s."
)

// Errors during parsing.
const (
	errDuplicateElement   = "Only %d element(s) with key %s permitted in scope %s."
	errDuplicateField     = "Only %d field(s) with key %s permitted in scope %s."
	errInvalidFieldValue  = "Value %s does not match regex %s for field %s."
	errInvalidToken       = "Invalid token %s in %s."
	errInsufficientFields = "Insufficient fields with key %s in element %s."
)
