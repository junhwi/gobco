# GOBCO - Golang Branch Coverage

Branch coverage measurement tool for golang.

## Install and Usage
```
$ go get github.com/junhwi/gobco/...
```

Add `ReportCoverage` to your `TestMain` function, for example:
```go
package package_name

import (
	"os"
	"testing"

	"github.com/junhwi/gobco"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	gobco.ReportCoverage()
	gobco.ReportProfile("c.out")
	os.Exit(retCode)
}
```

Then run `go test` with `-toolexec` flag:
```
$ go test -cover -toolexec 'gobco-tool'
--- FAIL: TestFoo (0.00s)
  foo_test.go:16: wrong
FAIL
bar.go 1 / 2
foo.go 5 / 6
exit status 1
FAIL  gobco/sample 0.008s
```

Generate an HTML report from profile:
```
$ gobco -html=c.out -o gobco.html
```
