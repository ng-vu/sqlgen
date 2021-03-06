package main

import (
	"bytes"
	"fmt"
	"go/types"
	"io"

	gen "github.com/ng-vu/sqlgen/gen/sqlgen"
)

type Adapter struct {
	pkg *types.Package
	b   bytes.Buffer
	n   int

	imports  map[string]string
	typeStrs map[types.Type]string
}

var _ gen.Interface = &Adapter{}

func NewAdapter() *Adapter {
	return &Adapter{
		imports:  make(map[string]string),
		typeStrs: make(map[types.Type]string),
	}
}

func (a *Adapter) WriteTo(w io.Writer) {
	fmt.Fprintf(w, "// Code generated by sqlgen DO NOT EDIT.\n\n")
	fmt.Fprintf(w, "package %v\n\n", a.pkg.Name())
	fmt.Fprintf(w, "import (\n")
	for path, name := range a.imports {
		if name != "" {
			fmt.Fprint(w, name, " ")
			fmt.Fprint(w, `"`, path, "\"\n")
		}
	}
	fmt.Fprintf(w, ")\n\n")
	w.Write(a.b.Bytes())
}

func (a *Adapter) P(format string, args ...interface{}) {
	for i := 0; i < a.n; i++ {
		a.b.WriteByte('\t')
	}
	fmt.Fprintf(&a.b, format, args...)
	a.b.WriteByte('\n')
}

func (a *Adapter) In() {
	a.n++
}

func (a *Adapter) Out() {
	if a.n <= 0 {
		panic("Negative")
	}
	a.n--
}

func (a *Adapter) NewImport(name, path string) func() string {
	a.imports[path] = name
	return func() string { return "" }
}

func (a *Adapter) TypeString(t types.Type) string {
	s, ok := a.typeStrs[t]
	if !ok {
		s = types.TypeString(t, a.Qualify)
	}
	return s
}

func (a *Adapter) Qualify(pkg *types.Package) string {
	if pkg == a.pkg {
		return ""
	}
	if _, ok := a.imports[pkg.Path()]; !ok {
		a.imports[pkg.Path()] = pkg.Name()
	}
	return pkg.Name()
}
