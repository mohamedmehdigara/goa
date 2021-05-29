package tests

import (
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
)

type (
	// Validator validates the generated code.
	Validator interface {
		// Validate validates the generated code and prints any errors.
		Validate(t *testing.T, file *codegen.File)
	}

	// ValidatorFunc is the type of functions invoked to validate generated code.
	ValidatorFunc func(t *testing.T, genpkg string, root *expr.RootExpr)

	// SectionValidator validates a given section. It implements Validator.
	SectionValidator struct {
		name     string
		expected string
	}
)

// ValidateSection returns a new SectionValidator.
func ValidateSection(name, expected string) *SectionValidator {
	return &SectionValidator{name: name, expected: expected}
}

// Validate validates the actual generated code with the expected code.
func (s *SectionValidator) Validate(t *testing.T, file *codegen.File) {
	var st *codegen.SectionTemplate
	{
		st = codegen.SectionTemplateWithName(t, file.SectionTemplates, s.name)
		switch {
		case st == nil && s.expected != "":
			t.Fatalf("%s: section %q not found when expected", file.Path, s.name)
		case st == nil && s.expected == "":
			// section has no code as expected
			return
		}
	}

	code := codegen.SectionCode(t, st)
	if code != s.expected {
		t.Errorf("invalid section code in %s: %s\n\ngot:\n%s\nexpected:\n%s\n", file.Path, s.name, code, codegen.Diff(t, code, s.expected))
	}
}
