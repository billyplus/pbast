package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/billyplus/pbast"
)

type printer struct {
}

func Fprint(w io.Writer, n pbast.Node) {
	p := &printer{}
	p.Fprint(w, n)
}

func (p *printer) Fprint(w io.Writer, n pbast.Node) {
	switch n := n.(type) {
	case *pbast.File:
		p.printFile(w, n)
	case pbast.Syntax:
		p.printSyntax(w, n)
	case *pbast.Import:
		p.printImport(w, n)
	case pbast.Package:
		p.printPackage(w, n)
	case *pbast.Option:
		p.printOption(w, n)
	case *pbast.Message:
		p.printMessage(w, n)
	case *pbast.MessageField:
		p.printMessageField(w, n)
	case *pbast.OneOf:
		p.printOneOf(w, n)
	case *pbast.OneOfField:
		p.printOneOfField(w, n)
	case *pbast.FieldOption:
		p.printFieldOption(w, n)
	case *pbast.Enum:
		p.printEnum(w, n)
	case *pbast.EnumField:
		p.printEnumField(w, n)
	case *pbast.EnumValueOption:
		p.printEnumValueOption(w, n)
	case *pbast.Service:
		p.printService(w, n)
	case *pbast.RPC:
		p.printRPC(w, n)
	case *pbast.ReturnType:
		p.printReturnType(w, n)
	case pbast.Comment:
		p.printComment(w, n)
	}
}

func (p *printer) printFile(w io.Writer, f *pbast.File) {
	// Comment
	p.Fprint(w, f.Comment)
	// syntax
	p.Fprint(w, f.Syntax)
	// imports
	for _, i := range f.Imports {
		p.Fprint(w, i)
	}
	// packages
	if f.Package != "" {
		p.Fprint(w, f.Package)
	}
	// options
	for _, o := range f.Options {
		p.Fprint(w, o)
	}
	// messages
	for _, m := range f.Messages {
		fmt.Fprintln(w)
		p.Fprint(w, m)
	}
	// enums
	for _, e := range f.Enums {
		fmt.Fprintln(w)
		p.Fprint(w, e)
	}
	// services
	for _, s := range f.Services {
		fmt.Fprintln(w)
		p.Fprint(w, s)
	}
}

func (p *printer) printSyntax(w io.Writer, s pbast.Syntax) {
	fmt.Fprintf(w, "syntax = \"%s\";", s)
	fmt.Fprintln(w)
}

func (p *printer) printImport(w io.Writer, i *pbast.Import) {
	if i.Visibility == pbast.NotSpecified {
		fmt.Fprintf(w, "import \"%s\";", i.Name)
		fmt.Fprintln(w)
		return
	}
	fmt.Fprintf(w, "import %s \"%s\";", i.Visibility, i.Name)
	fmt.Fprintln(w)
}

func (p *printer) printPackage(w io.Writer, pkg pbast.Package) {
	fmt.Fprintf(w, "package %s;", pkg)
	fmt.Fprintln(w)
}

func (p *printer) printOption(w io.Writer, o *pbast.Option) {
	fmt.Fprintf(w, "%s = %s;", o.Name, o.Value)
	fmt.Fprintln(w)
}

func (p *printer) printMessage(w io.Writer, m *pbast.Message) {
	// comment
	p.Fprint(w, m.Comment)

	// name
	fmt.Fprintf(w, "message %s {", m.Name)
	fmt.Fprintln(w)

	indent := pbast.NewSpaceWriter(w, shift)
	// fields
	for _, f := range m.Fields {
		p.Fprint(indent, f)
	}
	// enums
	for _, e := range m.Enums {
		p.Fprint(indent, e)
	}
	// messages
	for _, m := range m.Messages {
		p.Fprint(indent, m)
	}
	// oneofs
	for _, o := range m.OneOfs {
		p.Fprint(indent, o)
	}

	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printMessageField(w io.Writer, f *pbast.MessageField) {
	// comment
	p.printComment(w, f.Comment)

	if f.Repeated {
		fmt.Fprintf(w, "repeated ")
	}
	fmt.Fprintf(w, "%s %s = %d", f.Type, f.Name, f.Index)

	if len(f.Options) > 0 {
		fmt.Fprint(w, " [")
		p.Fprint(w, f.Options[0])

		for _, f := range f.Options[1:] {
			fmt.Fprint(w, ", ")
			p.Fprint(w, f)
		}
		fmt.Fprint(w, "]")
	}
	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printOneOf(w io.Writer, o *pbast.OneOf) {
	// comment
	p.printComment(w, o.Comment)

	// name
	fmt.Fprintf(w, "oneof %s {", o.Name)
	fmt.Fprintln(w)

	indent := pbast.NewSpaceWriter(w, shift)
	// fields
	for _, f := range o.Fields {
		p.Fprint(indent, f)
	}

	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printOneOfField(w io.Writer, f *pbast.OneOfField) {
	// comment
	p.printComment(w, f.Comment)

	fmt.Fprintf(w, "%s %s = %d", f.Type, f.Name, f.Index)

	if len(f.Options) > 0 {
		fmt.Fprint(w, " [")
		p.Fprint(w, f.Options[0])

		for _, f := range f.Options[1:] {
			fmt.Fprint(w, ", ")
			p.Fprint(w, f)
		}
		fmt.Fprint(w, "]")
	}
	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printFieldOption(w io.Writer, o *pbast.FieldOption) {
	fmt.Fprintf(w, "%s = %s", o.Name, o.Value)
}

func (p *printer) printEnum(w io.Writer, e *pbast.Enum) {
	// comment
	p.Fprint(w, e.Comment)
	// name
	fmt.Fprintf(w, "enum %s {", e.Name)
	fmt.Fprintln(w)
	// options
	for _, o := range e.Options {
		p.Fprint(w, o)
	}
	// fields
	for _, f := range e.Fields {
		p.Fprint(pbast.NewSpaceWriter(w, shift), f)
	}
	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printEnumField(w io.Writer, f *pbast.EnumField) {
	fmt.Fprintf(w, "%s = %d", f.Name, f.Index)

	if len(f.Options) != 0 {
		fmt.Fprintf(w, " [")
		opts := []string{}
		for _, o := range f.Options {
			opts = append(opts, fmt.Sprintf("%s = %s", o.Name, o.Value))
		}
		fmt.Fprint(w, strings.Join(opts, ", "))
		fmt.Fprintf(w, "]")
	}

	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printEnumValueOption(w io.Writer, o *pbast.EnumValueOption) {
	fmt.Fprintf(w, "%s = %s", o.Name, o.Value)
}

func (p *printer) printService(w io.Writer, s *pbast.Service) {
	// comment
	p.Fprint(w, s.Comment)

	fmt.Fprintf(w, "service %s {\n", s.Name)

	indent := pbast.NewSpaceWriter(w, shift)
	// options
	for _, o := range s.Options {
		p.Fprint(indent, o)
	}
	// RPCs
	for _, r := range s.RPCs {
		p.Fprint(indent, r)
	}

	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printRPC(w io.Writer, r *pbast.RPC) {
	// comment
	p.Fprint(w, r.Comment)

	fmt.Fprintf(w, "rpc %s ", r.Name)
	p.Fprint(w, r.Input)
	fmt.Fprint(w, " returns ")
	p.Fprint(w, r.Output)
	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printReturnType(w io.Writer, i *pbast.ReturnType) {
	fmt.Fprint(w, "(")
	if i.Streamable {
		fmt.Fprint(w, "stream ")
	}
	fmt.Fprintf(w, "%s)", i.Name)
}

func (p *printer) printComment(w io.Writer, c pbast.Comment) {
	lines := make([]string, 0, len(c))
	for _, line := range c {
		lines = append(lines, "// "+line+"\n")
	}

	fmt.Fprint(w, strings.Join(lines, ""))
}

const shift = 2
