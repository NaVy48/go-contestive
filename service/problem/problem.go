package problem

import "encoding/xml"

type Problem struct {
	XMLName    xml.Name           `xml:"problem"`
	Revision   int                `xml:"revision,attr"`
	ShortName  string             `xml:"short-name,attr"`
	Url        string             `xml:"url,attr"`
	Names      []ProblemName      `xml:"names>name"`
	Statements []ProblemStatement `xml:"statements>statement"`
	Testset    ProblemTestset     `xml:"judging>testset"`
}

type ProblemName struct {
	Value string `xml:"value,attr"`
	Lang  string `xml:"language,attr"`
}

type ProblemStatement struct {
	Charset  string `xml:"charset,attr"`
	Language string `xml:"language,attr"`
	Mathjax  bool   `xml:"mathjax,attr"`
	Path     string `xml:"path,attr"`
	Type     string `xml:"type,attr"`
}

type ProblemTestset struct {
	Name              string        `xml:"name,attr"`
	TimeLimit         int           `xml:"time-limit"`
	MemoryLimit       int           `xml:"memory-limit"`
	TestCount         int           `xml:"test-count"`
	InputPathPattern  string        `xml:"input-path-pattern"`
	AnswerPathPattern string        `xml:"answer-path-pattern"`
	Tests             []ProblemTest `xml:"tests>test"`
}

type ProblemTest struct {
	Method string `xml:"method,attr"`
	Points string `xml:"points,attr"`
	Sample string `xml:"sample,attr"`
}

const langEn = "english"

func (p *Problem) Title() string {
	for _, n := range p.Names {
		if n.Lang == langEn {
			return n.Value
		}
	}

	if len(p.Names) == 0 {
		return ""
	}

	return p.Names[0].Value
}
