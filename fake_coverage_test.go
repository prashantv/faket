package faket

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/prashantv/faket/internal/want"
)

// TestTBCoverage ensures faket implements all methods of testing.TB.
// Since fakeTB embeds testing.TB so it'll always implement the testing.TB
// interface, but unimplemented methods will panic if called.
// This test ensures each method is explicitly implemented.
func TestTBCoverage(t *testing.T) {
	// We can't use reflection as fakeTB embeds testing.TB which causes reflect
	// to include testing.TB methods under fakeTB.
	ftMethods, err := findTypeMethods(".", reflect.TypeOf(fakeTB{}).Name())
	want.NoErr(t, err)

	ftSet := make(map[string]struct{})
	for _, m := range ftMethods {
		ftSet[m] = struct{}{}
	}

	tb := reflect.TypeOf((*testing.TB)(nil)).Elem()
	for i := 0; i < tb.NumMethod(); i++ {
		m := tb.Method(i)

		if m.PkgPath != "" {
			// unexported, ignore
			continue
		}

		mn := m.Name
		// TODO(prashant): Implement new methods added in 1.24 before release.
		if mn == "Chdir" || mn == "Context" {
			continue
		}
		if _, ok := ftSet[mn]; !ok {
			t.Errorf("faket missing testing.TB.%s", mn)
		}
	}
}

func findTypeMethods(path, typeName string) ([]string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}

	var methods []string
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if !ok || funcDecl.Recv == nil {
					continue
				}

				recv := funcDecl.Recv.List
				if len(recv) == 0 {
					continue
				}

				if receiverName(recv[0]) != typeName {
					continue
				}

				methods = append(methods, funcDecl.Name.Name)
			}
		}
	}

	return methods, nil
}

func receiverName(f *ast.Field) string {
	switch recv := f.Type.(type) {
	case *ast.Ident:
		return recv.Name
	case *ast.StarExpr:
		if ident, ok := recv.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	panic(fmt.Errorf("cannot parse type name from receiver: %#v", f.Type))
}
