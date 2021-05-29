package tests

import (
	"testing"

	codegentests "goa.design/goa/v3/codegen/tests"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

// ServerFile validates the HTTP server's server.go file.
func ServerFile(v codegentests.Validator) codegentests.ValidatorFunc {
	return func(t *testing.T, genpkg string, root *expr.RootExpr) {
		fs := httpcodegen.ServerFiles(genpkg, root)
		if len(fs) == 0 {
			t.Fatalf("no server files generated")
		}
		v.Validate(t, fs[0])
	}
}

// ServerEncodeDecodeFile validates the HTTP server's encode_decode.go file.
func ServerEncodeDecodeFile(v codegentests.Validator) codegentests.ValidatorFunc {
	return func(t *testing.T, genpkg string, root *expr.RootExpr) {
		fs := httpcodegen.ServerFiles(genpkg, root)
		if len(fs) == 0 {
			t.Fatalf("no server files generated")
		}
		v.Validate(t, fs[1])
	}
}

// ClientFile validates the HTTP client's client.go file.
func ClientFile(v codegentests.Validator) codegentests.ValidatorFunc {
	return func(t *testing.T, genpkg string, root *expr.RootExpr) {
		fs := httpcodegen.ClientFiles(genpkg, root)
		if len(fs) == 0 {
			t.Fatalf("no client files generated")
		}
		v.Validate(t, fs[0])
	}
}

// ClientEncodeDecodeFile validates the HTTP client's encode_decode.go file.
func ClientEncodeDecodeFile(v codegentests.Validator) codegentests.ValidatorFunc {
	return func(t *testing.T, genpkg string, root *expr.RootExpr) {
		fs := httpcodegen.ClientFiles(genpkg, root)
		if len(fs) == 0 {
			t.Fatalf("no client files generated")
		}
		v.Validate(t, fs[1])
	}
}
