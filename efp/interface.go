package efp

func PrototypeBytes(bytes []byte) (*protoElement, []string) {
	return nil, nil
}

func PrototypeString(prototype string) (*protoElement, []string) {
	return PrototypeBytes([]byte(prototype))
}

func PrototypeFile(prototype string) (*protoElement, []string) {
	return nil, nil
}

// ValidateBytes against the
func (p *protoElement) ValidateBytes(bytes []byte) (*element, []string) {
	return nil, nil
}

func (p *protoElement) ValidateString(prototype string) (*element, []string) {
	return ValidateBytes([]byte(prototype))
}

func (p *protoElement) ValidateFile(prototype string) (*element, []string) {
	return nil, nil
}
