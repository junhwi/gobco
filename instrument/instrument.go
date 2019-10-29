package instrument

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Branch struct {
	start token.Pos
	end   token.Pos
}

func getCounter(name string, idx int) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.Ident{
				Name: name,
			},
			Index: &ast.BasicLit{
				Kind:  token.INT,
				Value: fmt.Sprint(idx),
			},
		},
	}
}

func (ctx *context) newCounter(name string, cond ast.Expr) *ast.CallExpr {

	id := len(ctx.branches)
	ctx.branches = append(ctx.branches, Branch{cond.Pos(), cond.End()})

	return &ast.CallExpr{
		Fun: ast.NewIdent(fmt.Sprintf("%sCount", ctx.pkg)),
		Args: []ast.Expr{
			cond,
			getCounter(name+".TCount", id),
			getCounter(name+".FCount", id),
		},
	}
}

func (ctx *context) createVisitor(name string) func(ast.Node) bool {
	return func(n ast.Node) bool {
		switch x := n.(type) {
		// TODO: We need to handle ast.CaseCluase
		// TODO: We need to handle go routine related things
		// such as ast.SelectStmt
		case *ast.IfStmt:
			x.Cond = ctx.newCounter(name, x.Cond)
		case *ast.ForStmt:
			x.Cond = ctx.newCounter(name, x.Cond)
		case *ast.FuncDecl:
			return !ctx.self || x.Name.Name != "Count"
		}
		return true
	}
}

func (ctx *context) importPkg() {
	spec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%q", ctx.pkgPath),
		},
	}
	ctx.file.Imports = append(ctx.file.Imports, spec)
	ctx.file.Decls = append([]ast.Decl{&ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{spec},
	}}, ctx.file.Decls...)
}

type context struct {
	file     *ast.File
	pkg      string
	pkgPath  string
	self     bool
	branches []Branch
}

func Instrument(name string, fd *os.File, coverVar string) error {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return err
	}

	ctx := context{
		file: f,
		pkg:  "",
		pkgPath: "github.com/junhwi/gobco",
		self: f.Name.Name == "gobco",
	}
	if !ctx.self {
		ctx.importPkg()
		ctx.pkg = "gobco."
	}

	gobcoVar := "Gobco_" + coverVar
	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, ctx.createVisitor(gobcoVar))


	printer.Fprint(fd, fset, f)
	fmt.Fprintf(fd, `
var %s = struct {
	Count []uint32
	Pos []uint32
	NumStmt []uint16
} {}
`, coverVar)

	total := len(ctx.branches)
	fmt.Fprintf(fd, "var %s = %sCov {\n", gobcoVar, ctx.pkg)
	fmt.Fprintf(fd, "\tTCount: make([]int, %d),\n", total)
	fmt.Fprintf(fd, "\tFCount: make([]int, %d),\n", total)
	fmt.Fprintf(fd, "\tPos: []uint32 {\n")
	for i, b := range ctx.branches {
		start := fset.Position(b.start)
		end := fset.Position(b.end)
		fmt.Fprintf(fd, "\t\t%d, %d, %#x, // [%d]\n", start.Line, end.Line, (end.Column&0xFFFF)<<16|(start.Column&0xFFFF), i)
	}
	fmt.Fprintf(fd, "\t},\n")
	fmt.Fprintf(fd, "}\n")
	fmt.Fprintf(fd, `
func init() {
	%sRegisterCov(&%s, "%s")
}
`, ctx.pkg, gobcoVar, name)

	return nil
}
