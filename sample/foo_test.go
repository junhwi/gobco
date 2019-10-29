package sample

import (
	"github.com/junhwi/gobco"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	gobco.ReportCoverage()
	gobco.ReportProfile("gobco.cover.out")
	os.Exit(retCode)
}

func TestFoo(t *testing.T) {
	if !Foo(9) {
		t.Error("wrong")
	}
}
