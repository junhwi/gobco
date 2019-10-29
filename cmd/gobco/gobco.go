package main

import (
	"flag"
	"fmt"
	"github.com/junhwi/gobco/html"
	"github.com/junhwi/gobco/instrument"
	"os"
)

func getFd(out string) (*os.File, error) {
	if out == "" {
		return os.Stdout, nil
	} else {
		return os.Create(out)
	}
}

func main() {

	coverVar := flag.String("var", "Cov", "name of coverage variable to generate (default \"Cov\")")
	out := flag.String("o", "", "file for output; default: stdout")
	htmlOut := flag.String("html", "", "generate HTML representation of coverage profile")
	flag.Parse()

	fd, err := getFd(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	profile := *htmlOut
	if profile != "" {
		err = html.ToHtml(profile, fd)
	} else {
		file := flag.Arg(0)
		err = instrument.Instrument(file, fd, *coverVar)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
