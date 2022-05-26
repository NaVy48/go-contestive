package judge

import (
	_ "embed"
	"errors"
	"os"
	"testing"

	"github.com/matryer/is"
)

//go:embed testdata/example-a-plus-b-test-6.zip
var archive []byte

//go:embed testdata/OK.cpp
var program []byte

var problemID = 6
var revisionID = 4

//go:embed testdata/OK.cpp
var srcOK []byte

//go:embed testdata/CE.cpp
var srcCE []byte

//go:embed testdata/MLE.cpp
var srcMLE []byte

//go:embed testdata/RT.cpp
var srcRT []byte

//go:embed testdata/TLE.cpp
var srcTLE []byte

//go:embed testdata/WA.cpp
var srcWA []byte

func Test_submission_runTests(t *testing.T) {
	err := os.Chdir(testDirectory)
	if err != nil {
		t.Fatal("unable to open problems directory")
	}

	dir := "/tmp/judge"
	if !checkDir(dir) {
		err = unpackProblem(archive, dir)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	makeNewParams := func(src []byte) SubmissionParams {
		return SubmissionParams{
			JudgeId:      0,
			Src:          src,
			ProblemID:    int64(problemID),
			RevisionID:   int64(revisionID),
			Language:     "cpp",
			ProblemDir:   dir,
			TempDir:      dir,
			ShortCircuit: true,
		}
	}

	tests := []struct {
		name        string
		src         []byte
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "OK",
			src:         srcOK,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "CE",
			src:         srcCE,
			wantErr:     true,
			expectedErr: NewError(CompilationError),
		},
		{
			name:        "MLE",
			src:         srcMLE,
			wantErr:     true,
			expectedErr: NewError(SignalKill),
		},
		{
			name:        "RT",
			src:         srcRT,
			wantErr:     true,
			expectedErr: NewError(SignalKill),
		},
		{
			name:        "TLE",
			src:         srcTLE,
			wantErr:     true,
			expectedErr: NewError(TimeLimitExceeded),
		},
		{
			name:        "WA",
			src:         srcWA,
			wantErr:     true,
			expectedErr: NewError(WrongAnswer),
		},
	}
	if testing.Short() {
		tests = tests[:1]
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			s := newSubmission(makeNewParams(tt.src))
			err := s.RunTests()
			if tt.wantErr {
				is.True(err != nil) // expecting an error
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("submission.RunTests() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				is.NoErr(err)
			}
		})
	}
}
