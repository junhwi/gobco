package main

import (
	"flag"
	"fmt"
	"github.com/junhwi/gobco/instrument"
	"log"
	"os"
	"os/exec"
	"path"
)

func getFd(out string) (*os.File, error) {
	if out == "" {
		return os.Stdout, nil
	} else {
		return os.Create(out)
	}
}

func runGobco() {

	cmd := flag.NewFlagSet("gobco", flag.ExitOnError)
	// Register all flags same as go tool cover
	outPtr := cmd.String("o", "", "file for output; default: stdout")
	version := cmd.String("V", "", "print version and exit")
	cmd.String("mode", "", "coverage mode: set, count, atomic")
	coverVar := cmd.String("var", "Cov", "name of coverage variable to generate (default \"Cov\")")
	cmd.Parse(os.Args[2:])
	files := cmd.Args()

	if *version != "" {
		fmt.Println("cover version go1.13.1")
	} else {
		fd, err := getFd(*outPtr)
		err = instrument.Instrument(files[0], fd, *coverVar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("gobco: ")

	tool := os.Args[1]
	args := os.Args[2:]

	toolName := path.Base(tool)
	if toolName == "cover" {
		runGobco()
	} else {
		cmd := exec.Command(tool, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(0)
}
