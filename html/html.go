package html

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type Profile struct {
	FileName   string
	Conditions []Condition
	Src        template.HTML
}

type Condition struct {
	StartLine, StartCol   int
	EndLine, EndCol       int
	TrueCount, FalseCount int
}

func parseProfiles(fileName string) ([]*Profile, error) {
	pf, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer pf.Close()

	files := make(map[string]*Profile)
	buf := bufio.NewReader(pf)
	s := bufio.NewScanner(buf)
	for s.Scan() {
		line := s.Text()
		name, cond, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		p := files[name]
		if p == nil {
			p = &Profile{
				FileName: name,
			}
			files[name] = p
		}
		p.Conditions = append(p.Conditions, cond)
	}

	profiles := make([]*Profile, 0, len(files))
	for _, p := range files {
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func parseLine(line string) (fileName string, cond Condition, err error) {
	end := len(line)
	c := Condition{}
	c.FalseCount, end, err = getNext(line, end)
	if err != nil {
		return "", c, err
	}
	c.TrueCount, end, err = getNext(line, end)
	if err != nil {
		return "", c, err
	}
	c.EndCol, end, err = getNext(line, end)
	if err != nil {
		return "", c, err
	}
	c.EndLine, end, err = getNext(line, end)
	if err != nil {
		return "", c, err
	}
	c.StartCol, end, err = getNext(line, end)
	if err != nil {
		return "", c, err
	}
	c.StartLine, end, err = getNext(line, end)
	if err != nil {
		return "", c, err
	}
	name := line[:end]
	if name == "" {
		return "", c, errors.New("empty FileName")
	}
	return name, c, nil
}

func getNext(s string, end int) (value int, next int, err error) {

	for i := end - 1; i >= 0; i-- {
		if s[i] == ',' {
			value, err := strconv.Atoi(s[i+1 : end])
			if err != nil {
				return 0, 0, err
			}
			return value, i, nil
		}
	}
	return 0, 0, errors.New("couldn't find a , ")
}

type Boundary struct {
	Offset     int
	Start      bool
	TrueCount  int
	FalseCount int
}

type conditionByStart []Condition

func (c conditionByStart) Len() int      { return len(c) }
func (c conditionByStart) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c conditionByStart) Less(i, j int) bool {
	ci, cj := c[i], c[j]
	return ci.StartLine < cj.StartLine || (ci.StartLine == cj.StartLine && ci.StartCol < cj.StartCol)
}

func (p *Profile) boundaries(src []byte) (boundaries []Boundary) {

	sort.Sort(conditionByStart(p.Conditions))

	line, col := 1, 1
	for srcIdx, condIdx := 0, 0; srcIdx < len(src) && condIdx < len(p.Conditions); {
		c := p.Conditions[condIdx]
		if c.StartLine == line && c.StartCol == col {
			boundaries = append(boundaries, Boundary{
				Offset:     srcIdx,
				Start:      true,
				TrueCount:  c.TrueCount,
				FalseCount: c.FalseCount,
			})
		}
		if c.EndLine == line && c.EndCol == col || line > c.EndLine {
			boundaries = append(boundaries, Boundary{Offset: srcIdx, Start: false, TrueCount: 0, FalseCount: 0})
			condIdx++
			continue
		}
		if src[srcIdx] == '\n' {
			line++
			col = 0
		}
		col++
		srcIdx++
	}
	return boundaries
}

func getColor(b Boundary) (color string) {
	if t, f := b.TrueCount > 0, b.FalseCount > 0; t && f {
		return "#4CAF50"
	} else if t || f {
		return "#ffeb3b"
	} else {
		return "#F44336"

	}
}

func htmlGen(w io.Writer, src []byte, boundaries []Boundary) error {
	dst := bufio.NewWriter(w)
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
			b := boundaries[0]
			color := getColor(b)

			if b.Start {
				dst.WriteString(fmt.Sprintf("<span style=\"color: %s; font-weight: bold\">", color))
			} else {
				dst.WriteString("</span>")
			}
			boundaries = boundaries[1:]
		}
		switch b := src[i]; b {
		case '>':
			dst.WriteString("&gt;")
		case '<':
			dst.WriteString("&lt;")
		case '&':
			dst.WriteString("&amp;")
		case '\t':
			dst.WriteString("        ")
		default:
			dst.WriteByte(b)
		}
	}
	return dst.Flush()
}

// TODO: Move to cmd
func ToHtml(profile string, out *os.File) error {

	html, err := template.New("html").Parse(htmlTmpl)
	if err != nil {
		return err
	}

	profiles, err := parseProfiles(profile)
	if err != nil {
		return err
	}
	for _, p := range profiles {
		src, err := ioutil.ReadFile(p.FileName)
		if err != nil {
			panic(err)
		}
		boundaries := p.boundaries(src)
		var buf bytes.Buffer
		err = htmlGen(&buf, src, boundaries)
		if err != nil {
			panic(err)
		}
		p.Src = template.HTML(buf.String())
	}
	err = html.Execute(out, profiles)

	return err
}

const htmlTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
	<div>
	{{range .}}
		<h1>{{.FileName}}</h1>
		<pre>{{.Src}}</pre>
	{{end}}
	</div>
</body>
</html>
`
