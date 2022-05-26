package judgemanager

import (
	"contestive/entity"
	"contestive/judgeconnection"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type handler struct {
	jm                *manager
	id                int
	conn              judgeconnection.Conn
	judgingSubmission int64
	notify            chan struct{}
}

func (h *handler) Run() error {
	for {
		submission, err := h.jm.SubmissionRepository.GetPendingSubmission(context.Background())
		if err != nil {
			if errors.Is(err, entity.ErrNotFound(nil)) {
				h.wait()
				continue
			}
			log.Printf("Encounterd error %v\n Seeping for some time and trying again", err)
			time.Sleep(time.Duration(rand.Int() % 1000 * 1000000)) // sleep 1-1000ms
		}

		err = h.judge(submission)
		if err != nil {
			log.Printf("Judgeing ended with: %v", err)
		}
	}
}

func (h *handler) judge(sub entity.Submission) (err error) {
	h.judgingSubmission = sub.ID
	defer func() { h.judgingSubmission = 0 }()

	sub.Status = "judging"
	h.jm.SubmissionRepository.UpdateSubmission(context.Background(), &sub)

	defer func() {
		if err != nil {
			sub.Status = "pending"
			h.jm.SubmissionRepository.UpdateSubmission(context.Background(), &sub)
		}
	}()

	h.conn.Write(judgeconnection.SubmitRequest{ProblemID: sub.ProblemID, RevisionID: sub.ProblemRevID, SubmissionID: sub.ID, Src: []byte(sub.SourceCode)})

	m := h.conn.Read()
	submitAck, ok := m.(judgeconnection.SubmitAck)
	if !ok || !submitAck.OK {
		return fmt.Errorf("unexpected message %T were expecting %T", m, submitAck)
	}

	m = h.conn.Read()
	packReq, ok := m.(judgeconnection.ProblemPackageRequest)
	if ok {
		pack, err := h.jm.ProblemRepository.PackageArchiveByRevisionId(context.Background(), packReq.RevisionID)
		if err != nil {
			return err
		}

		h.conn.Write(judgeconnection.ProblemPackageResponse{ProblemID: packReq.ProblemID, RevisionID: packReq.RevisionID, Package: pack})
		m = h.conn.Read()
	}

	result, ok := m.(judgeconnection.JudgeResult)
	if ok {
		sub.Result = result.Status
		sub.Details = result.Message
		sub.Status = "done"
		return nil
	}

	return fmt.Errorf("unexpected message %v", m)
}

func (h *handler) wait() {
	select {
	case <-h.notify:
	case <-h.conn.Closed():
		panic("Conn closed")
	}
}

func (h *handler) Notify() {
	select {
	case h.notify <- struct{}{}:
	default:
	}
}
