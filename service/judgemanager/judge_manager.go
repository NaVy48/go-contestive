package judgemanager

import (
	"contestive/entity"
	"contestive/judgeconnection"
	"context"
	"fmt"
	"log"
	"sync"
)

type SubmissionRepository interface {
	GetPendingSubmission(ctx context.Context) (entity.Submission, error)
	UpdateSubmission(ctx context.Context, s *entity.Submission) error
}
type ProblemRepository interface {
	PackageArchiveByRevisionId(ctx context.Context, revisionID entity.ID) ([]byte, error)
}

// Manager iterface is used for other system parts to interact with judge manager
type Manager interface {
	Notify()
	RunJudge(c judgeconnection.Conn) error
}

type manager struct {
	Secrets              map[int]string
	SubmissionRepository SubmissionRepository
	ProblemRepository    ProblemRepository
	judges               []*handler
	judgesMu             sync.Mutex
	srv                  judgeconnection.Server
}

// NewJudgeManager creats new judge manager instance. Runs TCP server for judges to connect and exposes interface to interact with them
func NewJudgeManager(address string, secrets map[int]string, srepo SubmissionRepository, prepo ProblemRepository) Manager {
	jm := &manager{Secrets: secrets, SubmissionRepository: srepo, ProblemRepository: prepo}
	jm.srv = judgeconnection.NewJudgeServer(jm, address)

	go func() {
		for {
			jm.listen()
		}
	}()

	return jm
}

func (jm *manager) listen() {
	var err error
	defer func() {
		if pan := recover(); pan != nil {
			log.Printf("Paniced while judge server was listening: %v", pan)
		} else if err != nil {
			log.Printf("Erro encounter while judge server was listening: %v", err)
		}
	}()
	err = jm.srv.Listen()
}

func (jm *manager) Close() {
	defer func() {
		if pan := recover(); pan != nil {
			log.Printf("Paniced while closing judge server: %v", pan)
		}
	}()
	jm.srv.Close()
}

func (jm *manager) Notify() {
	jm.judgesMu.Lock()
	defer jm.judgesMu.Unlock()

	for _, v := range jm.judges {
		v.Notify()
	}
}

// RunJudge runs judge controller on contest system side. Should return only when connection should be closed
func (jm *manager) RunJudge(c judgeconnection.Conn) error {
	message := c.Read()
	authReq, ok := message.(judgeconnection.AuthRequest)
	if !ok {
		return fmt.Errorf("auth failed")
	}
	h := &handler{jm, authReq.JudgeID, c, 0, make(chan struct{})}
	err := jm.authenticateJudge(h, authReq.Secret)
	if err != nil {
		return err
	}
	defer jm.removeJudge(h)

	c.Write(judgeconnection.AuthResponse{OK: true})

	return h.Run()
}

func (jm *manager) authenticateJudge(h *handler, secret string) error {
	jm.judgesMu.Lock()
	defer jm.judgesMu.Unlock()
	for _, v := range jm.judges {
		if v.id == h.id {
			return fmt.Errorf("already running")
		}
	}

	if jm.Secrets[h.id] != secret {
		return fmt.Errorf("auth failed")
	}
	jm.judges = append(jm.judges, h)

	return nil
}

func (jm *manager) removeJudge(h *handler) {
	jm.judgesMu.Lock()
	defer jm.judgesMu.Unlock()

	for i, v := range jm.judges {
		if v.id == h.id {
			last := len(jm.judges) - 1
			jm.judges[i] = jm.judges[last]
			jm.judges[last] = nil
			jm.judges = jm.judges[:last]
		}
	}
}
