package problem

import (
	"archive/zip"
	"bytes"
	"contestive/entity"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
)

type problemArchive struct {
	FS      fs.FS
	Problem *Problem
}

var ErrFileNotFound = fmt.Errorf("file not found")
var ErrInvalidProblemXML = fmt.Errorf("invalid problem.xml")

func ProblemArchive(archive []byte) (problemArchive, error) {
	reader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return problemArchive{}, entity.ErrCustomWrapper("problem archive unable to open zip", err)
	}
	p := problemArchive{FS: reader}

	p.Problem, err = p.ProblemXML()
	if err != nil {
		return problemArchive{}, ErrInvalidProblemXML
	}

	return p, nil
}

func ProblemArchiveFS(fs fs.FS) (problemArchive, error) {
	p := problemArchive{FS: fs}

	var err error
	p.Problem, err = p.ProblemXML()
	if err != nil {
		return problemArchive{}, ErrInvalidProblemXML
	}

	return p, nil
}

func (p problemArchive) ProblemXML() (*Problem, error) {
	if p.Problem != nil {
		return p.Problem, nil
	}
	f, err := p.FS.Open("problem.xml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	problem := new(Problem)
	err = xml.NewDecoder(f).Decode(problem)
	if err != nil {
		return nil, err
	}

	return problem, nil
}

func (p problemArchive) StatementPdf() ([]byte, error) {
	for _, s := range p.Problem.Statements {
		if s.Language == "english" && s.Type == "application/pdf" {
			f, err := p.FS.Open(s.Path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			return io.ReadAll(f)
		}
	}

	return nil, ErrFileNotFound
}

func (p problemArchive) StatementHtml() ([]byte, error) {
	for _, s := range p.Problem.Statements {
		if s.Language == "english" && s.Type == "text/html" {
			f, err := p.FS.Open(s.Path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			return io.ReadAll(f)
		}
	}

	return nil, ErrFileNotFound
}
