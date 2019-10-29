package instrument

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"testing"

	"github.com/junhwi/gobco"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	gobco.ReportCoverage()
	gobco.ReportProfile("instrument.out")
	os.Exit(retCode)
}

func Test_visitor(t *testing.T) {

	fs := token.NewFileSet()
	src := "func() { if True {} else {} }"
	node, _ := parser.ParseExpr(src)

	ctx := context{
		pkg: "instrument.",
	}
	visitor := ctx.createVisitor("name")
	ast.Inspect(node, visitor)

	buf := new(bytes.Buffer)
	printer.Fprint(buf, fs, node)
	assert.Equal(t, `func() {
	if instrument.Count(True, &name.TCount[0], &name.FCount[0]) {
	} else {
	}
}`, buf.String())
}

func Test_getCounter(t *testing.T) {

	fs := token.NewFileSet()

	expr := getCounter("name", 1)

	buf := new(bytes.Buffer)
	printer.Fprint(buf, fs, &ast.ExprStmt{X: expr})
	assert.Equal(t, "&name[1]", buf.String())
}

func Test_newCounter(t *testing.T) {

	fs := token.NewFileSet()
	cond, _ := parser.ParseExpr("1 == 2")

	ctx := context{
		pkg: "instrument.",
	}
	expr := ctx.newCounter("name", cond)

	buf := new(bytes.Buffer)
	printer.Fprint(buf, fs, &ast.ExprStmt{X: expr})
	assert.Equal(t, "instrument.Count(1 == 2, &name.TCount[0], &name.FCount[0])", buf.String())
}
