package judge

import (
	"bytes"
	"contestive/service/problem"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type SubmissionParams struct {
	JudgeId      int
	Src          []byte
	ProblemID    int64
	RevisionID   int64
	Language     string
	ProblemDir   string
	TempDir      string
	ShortCircuit bool
}

type submission struct {
	// params
	SubmissionParams

	// created directiries and files
	solutionTempDir string
	solutionSrcName string
	solutionBinName string
	isolatedDir     string
	boxBinFileName  string

	// result
	tests []test
	err   error
}

func newSubmission(params SubmissionParams) *submission {
	return &submission{SubmissionParams: params}
}

func (s *submission) makeTempDir() {
	solutionDir, err := os.MkdirTemp(s.TempDir, "temp")
	if err != nil {
		panic(err)
	}
	s.solutionTempDir = solutionDir
}

func (s *submission) cleanupTempDir() {
	if s.solutionTempDir != "" {
		err := os.RemoveAll(s.solutionTempDir)
		if err != nil {
			panic(err)
		}

		s.solutionTempDir = ""
	}
}

func (s *submission) createSrcFile() {
	s.solutionSrcName = path.Join(s.solutionTempDir, "solution.cpp")
	file, err := os.Create(s.solutionSrcName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(file, bytes.NewReader(s.Src))
	if err != nil {
		panic(err)
	}
}

func (s *submission) compileSrc() error {
	s.createSrcFile()
	s.solutionBinName = path.Join(s.solutionTempDir, "solution")

	err := CppCompile(s.solutionSrcName, s.solutionBinName)
	if err != nil {
		s.err = NewErrorWrap(CompilationError, err)
		return s.err
	}

	return nil
}

func (s *submission) initIsolate() {
	judgeId := strconv.Itoa(s.JudgeId)
	exec.Command("isolate", "-b", judgeId, "--cleanup").Run() // if by some chance there is hanging session
	initCmd := exec.Command("isolate", "-b", judgeId, "--init")
	initCmd.Stderr = os.Stderr
	output, err := initCmd.Output()
	log.Println(string(output))
	if err != nil {
		panic(err)
	}

	s.isolatedDir = path.Clean(strings.TrimSpace(string(output)))
}

func (s *submission) cleanupIsolate() {
	judgeId := strconv.Itoa(s.JudgeId)
	exec.Command("isolate", "-b", judgeId, "--cleanup").Run() // if by some chance there is hanging session
}

func (s *submission) prepareBox() {
	s.boxBinFileName = path.Join(s.isolatedDir, "box", "solution")
	err := copyFile(s.solutionBinName, s.boxBinFileName, 0770)
	if err != nil {
		panic(err)
	}
}

func (s *submission) runTests() error {
	s.makeTempDir()
	defer s.cleanupTempDir()

	err := s.compileSrc()
	if err != nil {
		return err
	}

	s.initIsolate()
	defer s.cleanupIsolate()

	s.prepareBox()

	pa, err := problem.ProblemArchiveFS(os.DirFS(s.ProblemDir))
	if err != nil {
		panic(err)
	}
	testset := pa.Problem.Testset
	inputFileFormat := path.Join(s.ProblemDir, testset.InputPathPattern)
	answerFileFormat := path.Join(s.ProblemDir, testset.AnswerPathPattern)
	outputFileFormat := path.Join(s.solutionTempDir, "%02d.out")

	RunReport := path.Join(s.solutionTempDir, "%02d.runr")
	CheckReprot := path.Join(s.solutionTempDir, "%02d.checkr")

	testCount := testset.TestCount
	s.tests = make([]test, 0, testCount)

	for i := 1; i <= testCount; i++ {
		t := test{
			index:            i,
			inFileName:       fmt.Sprintf(inputFileFormat, i),
			outFileName:      fmt.Sprintf(outputFileFormat, i),
			ansFileName:      fmt.Sprintf(answerFileFormat, i),
			runReportFName:   fmt.Sprintf(RunReport, i),
			checkReprotFName: fmt.Sprintf(CheckReprot, i),
			t:                testset,
		}

		s.tests = append(s.tests, t)

		t.runTest()
		t.runChecker(path.Join(s.ProblemDir, "files", "check"))

		if s.ShortCircuit && t.err != nil {
			s.err = NewFailedTestError(t)
			return s.err
		}
	}

	return nil
}

func (s *submission) RunTests() (err error) {
	defer func() {
		if pan := recover(); pan != nil {
			switch p := pan.(type) {
			case error:
				err = NewErrorWrap(JudgeError, p)
			case string:
				err = NewErrorWrap(JudgeError, fmt.Errorf("%s", p))
			default:
				log.Panicln(p)
				err = NewError(JudgeError)
			}
		}
	}()

	return s.runTests()
}
