package gobco

import (
	"fmt"
	"os"
	"path"
)

type Cov struct {
	TCount []int
	FCount []int
	Pos    []uint32
}

var cov = make(map[string]*Cov)

func Count(cond bool, true *int, false *int) bool {
	if cond {
		*true++
	} else {
		*false++
	}
	return cond
}

func RegisterCov(c *Cov, fileName string) {
	cov[fileName] = c
}

func ReportCoverage() {

	for k, v := range cov {
		covered, total := calculateCoverage(v)
		fileName := path.Base(k)
		fmt.Println(fileName, covered, "/", total)
	}
}

func ReportProfile(file string) error {

	f, err := os.Create(file)
	if err != nil {
		defer func() { f.Close() }()
		return err
	}

	for file, v := range cov {
		for i := range v.TCount {
			startLine := v.Pos[3*i+0]
			startCol := uint16(v.Pos[3*i+2])
			endLine := v.Pos[3*i+1]
			endCol := uint16(v.Pos[3*i+2] >> 16)
			fmt.Fprintf(f, "%s,%d,%d,%d,%d,%d,%d\n", file,
				startLine, startCol,
				endLine, endCol,
				v.TCount[i], v.FCount[i])
		}
	}

	return nil
}

func calculateCoverage(cov *Cov) (int, int) {
	cnt := 0
	for _, c := range cov.TCount {
		if c > 0 {
			cnt++
		}
	}
	for _, c := range cov.FCount {
		if c > 0 {
			cnt++
		}
	}
	return cnt, len(cov.TCount) * 2
}
