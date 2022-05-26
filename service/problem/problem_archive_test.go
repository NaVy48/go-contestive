package problem

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
)

//go:embed testdata/example-a-plus-b-test-6.zip
var archive []byte

var pa problemArchive

func TestMain(m *testing.M) {
	var err error
	pa, err = ProblemArchive(archive)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func Test_ProblemArchive_ProblemXML(t *testing.T) {
	is := is.New(t)

	problem, err := pa.ProblemXML()
	is.NoErr(err)
	fmt.Println(problem)
}
