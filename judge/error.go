package judge

import (
	"fmt"
)

type ErrorStatus string

const (
	RunTimeError      ErrorStatus = "RE"
	SignalKill        ErrorStatus = "SG"
	TimeLimitExceeded ErrorStatus = "TO"
	SandBoxError      ErrorStatus = "XX"
	CompilationError  ErrorStatus = "CE"
	WrongAnswer       ErrorStatus = "WA"
	JudgeError        ErrorStatus = "JE"
	UnknownError      ErrorStatus = "UN"
)

func (status ErrorStatus) String() string {
	switch status {
	case RunTimeError:
		return "run time error"
	case SignalKill:
		return "signal kill error"
	case TimeLimitExceeded:
		return "time limit exceeded"
	case SandBoxError:
		return "sand box error"
	case CompilationError:
		return "compilation error"
	case WrongAnswer:
		return "Wrong answer"
	case JudgeError:
		return "Internal judge error"
	default:
		return "unknown error"
	}
}

type judgeError struct {
	status      ErrorStatus
	internalErr error
}

func (err judgeError) Error() string {
	if err.internalErr != nil {
		return fmt.Sprintf("%s: %s", err.status, err.internalErr.Error())
	}
	return err.status.String()
}

func (err judgeError) Unwrap() error {
	return err.internalErr
}

func (err judgeError) Is(target error) bool {
	t, ok := target.(judgeError)
	return ok && t.status == err.status
}

func (err judgeError) As(target interface{}) bool {
	t, ok := target.(*judgeError)
	if ok && t.status == err.status {
		*t = err
		return true
	}
	return false
}

func NewError(status ErrorStatus) error {
	return judgeError{status, nil}
}

func NewErrorWrap(status ErrorStatus, innerErr error) error {
	return judgeError{status, innerErr}
}

type failedTestError struct {
	testIndex  int
	innerError error
	status     ErrorStatus
}

func NewFailedTestError(t test) error {
	if t.err == nil {
		return nil
	}

	err := failedTestError{t.index, t.err, UnknownError}
	je, ok := t.err.(judgeError)
	if ok {
		err.status = je.status
	}
	return err
}

func (err failedTestError) Error() string {
	if err.innerError != nil {
		return fmt.Sprintf("Failed test %d with %v: %s", err.testIndex, string(err.status), err.innerError.Error())
	}
	return fmt.Sprintf("Failed test %d with %v", err.testIndex, string(err.status))
}

func (err failedTestError) Unwrap() error {
	return err.innerError
}
