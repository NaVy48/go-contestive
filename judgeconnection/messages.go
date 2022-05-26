package judgeconnection

type DTO struct {
	AuthRequest            *AuthRequest
	AuthResponse           *AuthResponse
	SubmitRequest          *SubmitRequest
	SubmitAck              *SubmitAck
	ProblemPackageRequest  *ProblemPackageRequest
	ProblemPackageResponse *ProblemPackageResponse
	JudgeResult            *JudgeResult
}

type Message interface {
	isMessage()
}

type AuthRequest struct {
	JudgeID int
	Secret  string
}

func (AuthRequest) isMessage() {}

type AuthResponse struct {
	OK bool
}

func (AuthResponse) isMessage() {}

type SubmitRequest struct {
	ProblemID    int64
	RevisionID   int64
	SubmissionID int64
	Src          []byte
}

func (SubmitRequest) isMessage() {}

type SubmitAck struct {
	OK bool
}

func (SubmitAck) isMessage() {}

type ProblemPackageRequest struct {
	ProblemID  int64
	RevisionID int64
}

func (ProblemPackageRequest) isMessage() {}

type ProblemPackageResponse struct {
	ProblemID  int64
	RevisionID int64
	Package    []byte
}

func (ProblemPackageResponse) isMessage() {}

type JudgeResult struct {
	ProblemID    int
	RevisionID   int
	SubmissionID int
	Status       string
	Message      string
	Test         []TestResult
	Error        string
}

func (JudgeResult) isMessage() {}

type TestResult struct {
	Index   int
	Status  string
	Message string
}
