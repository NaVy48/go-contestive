package judge

import (
	"contestive/service/problem"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	index            int
	inFileName       string
	outFileName      string
	ansFileName      string
	runReportFName   string
	checkReprotFName string
	t                problem.ProblemTestset
	// results
	time       float64
	wallTime   float64
	memoryUsed int
	exitcode   int
	err        error
}

func (t *test) runTest() {
	inFile, err := os.Open(t.inFileName)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	outFile, err := os.Create(t.outFileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	// "isolate -t 1s -w 10s --run -- solution"
	cmd := exec.Command("isolate", "-M", t.runReportFName, "-t", fmt.Sprintf("%f", float32(t.t.TimeLimit)/1000), "-w", fmt.Sprintf("%f", float32(t.t.TimeLimit)/1000*3), "-m", fmt.Sprintf("%d", t.t.MemoryLimit>>10), "--run", "--", "solution")
	cmd.Stdin = inFile
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr
	cmd.Run()

	t.checkTestReport()
}

func (t *test) checkTestReport() {
	rf, err := os.Open(t.runReportFName)
	if err != nil {
		panic(err)
	}
	defer rf.Close()

	report := make(map[string]string)
	temp := ""
	for err == nil {
		_, err = fmt.Fscanln(rf, &temp)
		if ind := strings.Index(temp, ":"); ind > 0 {
			key := temp[:ind]
			val := temp[ind+1:]
			report[key] = val
		}
	}

	fmt.Printf("\n\n#####\n %s \n#########\n\n", report["message"])

	if time, ok := report["time"]; ok {
		t.time, _ = strconv.ParseFloat(time, 64)
	}

	if timeWall, ok := report["time-wall"]; ok {
		t.wallTime, _ = strconv.ParseFloat(timeWall, 64)
	}

	if maxrss, ok := report["max-rss"]; ok {
		t.memoryUsed, _ = strconv.Atoi(maxrss)
		t.memoryUsed = t.memoryUsed << 10 // convert from kB to bytes
	}

	if exitcode, ok := report["exitcode"]; ok {
		t.exitcode, err = strconv.Atoi(exitcode)
		if err != nil {
			t.exitcode = 11
		}
	} else {
		t.exitcode = 1
	}

	if t.exitcode != 0 {
		status := report["status"]
		switch status {
		case string(RunTimeError):
			t.err = NewError(RunTimeError)
		case string(SignalKill):
			t.err = NewError(SignalKill)
		case string(TimeLimitExceeded):
			t.err = NewError(TimeLimitExceeded)
		case string(SandBoxError):
			t.err = NewError(SandBoxError)
		default:
			t.err = NewError(SandBoxError)
		}
	}
}

func (t *test) runChecker(checkerFileName string) {
	// cmd := exec.Command(path.Join(dir, "files", "check"), t.inFileName, t.outFileName, t.ansFileName, t.checkReprotFName)

	if t.exitcode != 0 {
		return
	}

	cmd := exec.Command(checkerFileName, t.inFileName, t.outFileName, t.ansFileName, t.checkReprotFName)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	fmt.Printf("starting checker...\n")
	err := cmd.Run()
	if err != nil {
		t.err = NewError(WrongAnswer)
	}

	fmt.Printf("checker Finished.\n")
}
