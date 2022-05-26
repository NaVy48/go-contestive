package judge

import (
	"contestive/judgeconnection"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

type Judge struct {
	connection    judgeconnection.Conn
	testDirectory string
	judgeId       int
}

func RunJudge(connection judgeconnection.Conn, testDirectory string, judgeId int) error {
	j := Judge{connection, testDirectory, judgeId}

	return j.Run()
}

func (j *Judge) problemDir(problemID int64, revisionID int64) string {
	return path.Join(j.testDirectory, strconv.Itoa(int(problemID)), strconv.Itoa(int(revisionID)))
}

func (j *Judge) Run() (err error) {
	defer func() {
		if pan := recover(); pan != nil {
			log.Printf("Recovering from judge panic: %v", pan)
			switch p := pan.(type) {
			case error:
				err = fmt.Errorf("encounter a panic: %w", p)
			default:
				err = fmt.Errorf("encounter a panic: %v", p)

			}
		}
	}()

	err = os.Chdir(j.testDirectory)
	if err != nil {
		return fmt.Errorf("unable to open problems directory %s", j.testDirectory)
	}

	for {
		val := j.connection.Read()
		submission, ok := val.(judgeconnection.SubmitRequest)
		if !ok {
			return fmt.Errorf("unexpected message type %T, expected judgeconnection.SubmitRequest", submission)
		}

		err = j.HandleSubmission(submission)
		if err != nil {
			return err
		}
	}
}

func (j *Judge) HandleSubmission(submission judgeconnection.SubmitRequest) error {
	dir := j.problemDir(submission.ProblemID, submission.RevisionID)
	if !checkDir(dir) {
		j.connection.Write(judgeconnection.ProblemPackageRequest{
			ProblemID:  submission.ProblemID,
			RevisionID: submission.RevisionID,
		})

		val := j.connection.Read()
		packageResponse, ok := val.(judgeconnection.ProblemPackageResponse)
		if !ok {
			return fmt.Errorf("unexpected message type %T, expected judgeconnection.SubmitRequest", submission)
		}

		err := unpackProblem(packageResponse.Package, dir)
		if err != nil {
			return err
		}
	}

	s := newSubmission(SubmissionParams{
		JudgeId:      j.judgeId,
		Src:          submission.Src,
		ProblemID:    submission.ProblemID,
		RevisionID:   submission.RevisionID,
		Language:     "cpp",
		ProblemDir:   j.problemDir(submission.ProblemID, submission.RevisionID),
		TempDir:      j.problemDir(submission.ProblemID, submission.RevisionID),
		ShortCircuit: true,
	})
	err := s.RunTests()
	if err != nil {
		return err
	}

	return nil
}
