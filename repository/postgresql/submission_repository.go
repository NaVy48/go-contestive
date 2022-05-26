package postgresql

import (
	"contestive/entity"
	"context"
	"fmt"
	"log"
	"strings"
)

type SubmissionRepository Store

const submissionTableName = `submission`

var submissionFieldNames = []string{
	"id",
	"createdat",
	"updatedat",
	"problemid",
	"problemrevid",
	"contestid",
	"authorid",
	"language",
	"sourcecode",
	"status",
	"result",
	"details",
}

var submissionFieldsQuery = strings.Join(submissionFieldNames, ", ")

func (s SubmissionRepository) Create(ctx context.Context, submission *entity.Submission) error {
	rows, err := s.db.NamedQueryContext(ctx, `
	INSERT INTO `+submissionTableName+`
				( problemid,  problemrevid,  contestid,  authorid,  sourcecode,  language,  status)
	VALUES(:problemid, :problemrevid, :contestid, :authorid, :sourcecode, :language, :status)
	RETURNING `+submissionFieldsQuery,
		submission)
	if err != nil {
		return handleError(err)
	}
	defer rows.Close()
	rows.Next()

	return handleError(rows.StructScan(submission))
}

func (s SubmissionRepository) GetByID(ctx context.Context, id entity.ID, userID entity.ID) (entity.Submission, error) {
	var submission entity.Submission
	var err error

	if userID != 0 {
		err = s.db.GetContext(ctx, &submission, `
		SELECT `+submissionFieldsQuery+`
		FROM `+submissionTableName+`
		WHERE authorid = $1 AND id = $2
		LIMIT 1;`,
			userID, id,
		)
	} else {
		err = s.db.GetContext(ctx, &submission, `
		SELECT `+submissionFieldsQuery+`
		FROM `+submissionTableName+`
		WHERE id = $1
		LIMIT 1;`,
			id,
		)
	}
	if err != nil {
		return entity.Submission{}, handleError(err)
	}

	return submission, handleError(err)
}

// ListAll gets submissions acording to options from all active submission list
func (s SubmissionRepository) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Submission, int, error) {
	var submissions []struct {
		entity.Submission
		Total int
	}

	qo := SqlListOptions(options, submissionFieldNames)

	query := s.db.Rebind(fmt.Sprintf(`SELECT %s, count(*) OVER() AS total FROM %s WHERE %s %s %s`,
		submissionFieldsQuery,
		submissionTableName,
		qo.filterSql(),
		qo.orderSql(),
		qo.rangeSql(),
	))

	err := s.db.SelectContext(ctx, &submissions, query, qo.filterVars...)
	if err != nil {
		return nil, 0, handleError(err)
	}

	res := make([]entity.Submission, len(submissions))
	for i, v := range submissions {
		res[i] = v.Submission
	}
	total := 0
	if len(submissions) > 0 {
		total = submissions[0].Total
	}

	return res, total, nil
}

// ListAllForUsers gets submissions acording to options from all active submission that belongs to the user
func (s SubmissionRepository) ListAllForUsers(ctx context.Context, options entity.ListOptions, userid entity.ID) ([]entity.Submission, int, error) {
	var submissions []struct {
		entity.Submission
		Total int
	}

	qo := SqlListOptions(options, submissionFieldNames)

	query := s.db.Rebind(fmt.Sprintf(`SELECT %s, count(*) OVER() AS total FROM %s WHERE authorid=%d AND %s %s %s`,
		submissionFieldsQuery,
		submissionTableName,
		userid,
		qo.filterSql(),
		qo.orderSql(),
		qo.rangeSql(),
	))

	log.Printf("%s \n vars: %v", query, qo.filterVars)

	err := s.db.SelectContext(ctx, &submissions, query, qo.filterVars...)

	if err != nil {
		return nil, 0, handleError(err)
	}

	res := make([]entity.Submission, len(submissions))
	for i, v := range submissions {
		res[i] = v.Submission
	}
	total := 0
	if len(submissions) > 0 {
		total = submissions[0].Total
	}

	return res, total, nil
}

func (s SubmissionRepository) UpdateSubmission(ctx context.Context, submission *entity.Submission) error {
	rows, err := s.db.NamedQueryContext(ctx, `
	UPDATE `+submissionTableName+`
	SET ( status,  result,  details) =
      (:status, :result, :details)
	WHERE id = :id
	RETURNING `+submissionFieldsQuery,
		submission)
	if err != nil {
		return handleError(err)
	}
	defer rows.Close()
	rows.Next()

	return handleError(rows.StructScan(submission))
}

func (s SubmissionRepository) GetPendingSubmission(ctx context.Context) (entity.Submission, error) {
	options := entity.ListOptions{
		Range: entity.RangeParam{From: 0, To: 1},
		Sort:  entity.SortParam{Field: "createdat", Asc: true},
		Filter: []entity.FilterParam{{
			Field:        "status",
			FilterValues: []interface{}{"pending"},
		}},
	}

	qo := SqlListOptions(options, submissionFieldNames)

	query := s.db.Rebind(fmt.Sprintf(`SELECT %s FROM %s WHERE %s %s %s`,
		submissionFieldsQuery,
		submissionTableName,
		qo.filterSql(),
		qo.orderSql(),
		qo.rangeSql(),
	))

	var submission entity.Submission
	err := s.db.GetContext(ctx, &submission, query, qo.filterVars...)
	if err != nil {
		return submission, handleError(err)
	}
	return submission, nil
}
