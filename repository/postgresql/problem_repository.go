package postgresql

import (
	"contestive/entity"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ProblemRepository Store

const problemTableName = `"problem"`
const revisionTableName = `"problem_revision"`

var problemFieldNames = []string{
	"id",
	"createdat",
	"updatedat",
	"authorid",
	"name",
	"externalurl",
}

var allowedProblemFilters = append(problemFieldNames, "cp.contestid")

var problemFieldsQuery = strings.Join(problemFieldNames, ", ")

var revisionFieldNames = []string{
	"id",
	"createdat",
	"updatedat",
	"authorid",
	"problemid",
	"revision",
	"title",
	"memorylimit",
	"timelimit",
	// "statementhtml",
	// "statementpdf",
	// "packagearchive",
	"outdated",
}
var revisionFieldsQuery = strings.Join(revisionFieldNames, ", ")

func (s ProblemRepository) listAllQuery(qo *sqlListOptions, userID entity.ID) string {
	if userID == 0 {
		return s.db.Rebind(fmt.Sprintf(`SELECT
	p.id id,
	p.createdat,
	p.updatedat,
	p.authorid,
	p.name,
	p.externalurl,
	rev.id revid,
	rev.createdat revcreatedat,
	rev.updatedat revupdatedat,
	rev.authorid revauthorid,
	rev.problemid,
	rev.revision,
	rev.title,
	rev.memorylimit,
	rev.timelimit,
	rev.outdated,
	count(*) OVER() AS total
	FROM %s p
	INNER JOIN %s rev ON rev.problemid = p.id AND rev.outdated = FALSE
	LEFT JOIN contest_problem cp ON cp.problemid = p.id
	WHERE %s %s %s`,
			problemTableName,
			revisionTableName,
			qo.filterSql(),
			qo.orderSql(),
			qo.rangeSql(),
		))
	} else {
		return s.db.Rebind(fmt.Sprintf(`SELECT
		p.id id,
		p.createdat,
		p.updatedat,
		p.authorid,
		p.name,
		p.externalurl,
		rev.id revid,
		rev.createdat revcreatedat,
		rev.updatedat revupdatedat,
		rev.authorid revauthorid,
		rev.problemid,
		rev.revision,
		rev.title,
		rev.memorylimit,
		rev.timelimit,
		rev.outdated,
		count(*) OVER() AS total
	FROM %s p
	INNER JOIN %s rev
		ON rev.problemid = p.id AND rev.outdated = FALSE
	INNER JOIN contest_problem cp
		ON cp.problemid = p.id
	INNER JOIN contest c
		ON cp.contestid = c.id AND c.startTime < NOW() AND c.endTime > NOW()
	INNER JOIN contest_user cu
		ON cu.contestid = c.id AND cu.userid = %d
	WHERE %s %s %s`,
			problemTableName,
			revisionTableName,
			userID,
			qo.filterSql(),
			qo.orderSql(),
			qo.rangeSql(),
		))

	}
}

func (s ProblemRepository) ListAll(ctx context.Context, options entity.ListOptions, userID entity.ID) ([]entity.Problem, int, error) {
	var problems []struct {
		entity.Problem
		entity.ProblemRevision
		RevID        int64
		RevCreatedAt time.Time
		RevUpdatedAt time.Time
		RevAuthorID  int64
		Total        int
	}

	fnames := make([]string, len(problemFieldNames))
	copy(fnames, problemFieldNames)
	fnames[0] = "p.id"
	for i := range options.Filter {
		if options.Filter[i].Field == "id" {
			options.Filter[i].Field = "p.id"
		}
	}
	qo := SqlListOptions(options, allowedProblemFilters)
	query := s.listAllQuery(&qo, userID)
	fmt.Println(query)

	err := s.db.SelectContext(ctx, &problems, query, qo.filterVars...)

	if err != nil {
		return nil, 0, handleError(err)
	}

	res := make([]entity.Problem, len(problems))
	for i, v := range problems {
		rev := v.ProblemRevision
		rev.ID = v.RevID
		rev.CreatedAt = v.RevCreatedAt
		rev.UpdatedAt = v.RevUpdatedAt
		rev.AuthorID = v.RevAuthorID
		v.Problem.Revisions = []entity.ProblemRevision{rev}
		res[i] = v.Problem
	}
	total := 0
	if len(problems) > 0 {
		total = problems[0].Total
	}

	return res, total, nil
}

func (s ProblemRepository) getByIDQuery(userID entity.ID) string {
	if userID == 0 {
		return fmt.Sprintf(`SELECT
		p.id,
		p.createdat,
		p.updatedat,
		p.authorid,
		p.name,
		p.externalurl,
		rev.id revid,
		rev.createdat revcreatedat,
		rev.updatedat revupdatedat,
		rev.authorid revauthorid,
		rev.problemid,
		rev.revision,
		rev.title,
		rev.memorylimit,
		rev.timelimit,
		rev.statementhtml,
		rev.outdated,
		count(*) OVER() AS total FROM %s p LEFT JOIN %s rev ON rev.problemid = p.id WHERE ProblemID = $1`,
			problemTableName,
			revisionTableName,
		)
	} else {
		return fmt.Sprintf(`SELECT
		p.id,
		p.createdat,
		p.updatedat,
		p.authorid,
		p.name,
		p.externalurl,
		rev.id revid,
		rev.createdat revcreatedat,
		rev.updatedat revupdatedat,
		rev.authorid revauthorid,
		rev.problemid,
		rev.revision,
		rev.title,
		rev.memorylimit,
		rev.timelimit,
		rev.statementhtml,
		rev.outdated,
		count(*) OVER() AS total
		FROM %s p
		INNER JOIN %s rev ON rev.problemid = p.id
		INNER JOIN contest_problem cp
			ON cp.problemid = p.id
		INNER JOIN contest c
			ON cp.contestid = c.id AND c.startTime < NOW() AND c.endTime > NOW()
		INNER JOIN contest_user cu
			ON cu.contestid = c.id AND cu.userid = %d
		WHERE p.id = $1`,
			problemTableName,
			revisionTableName,
			userID,
		)
	}
}
func (s ProblemRepository) GetByID(ctx context.Context, problemID entity.ID, userID entity.ID) (entity.Problem, error) {
	var problems []struct {
		entity.Problem
		entity.ProblemRevision
		RevID        int64
		RevCreatedAt time.Time
		RevUpdatedAt time.Time
		RevAuthorID  int64
		Total        int
	}

	err := s.db.SelectContext(ctx, &problems, s.getByIDQuery(userID), problemID)
	if err != nil {
		return entity.Problem{}, handleError(err)
	}

	if len(problems) == 0 {
		return entity.Problem{}, entity.ErrNotFound(nil)

	}

	res := problems[0].Problem

	revisions := make([]entity.ProblemRevision, len(problems))
	for i, v := range problems {
		revisions[i] = v.ProblemRevision
		revisions[i].ID = v.RevID
		revisions[i].CreatedAt = v.RevCreatedAt
		revisions[i].UpdatedAt = v.RevUpdatedAt
		revisions[i].AuthorID = v.RevAuthorID
	}

	res.Revisions = revisions

	return res, nil
}

func (s ProblemRepository) createRevision(tx *sqlx.Tx, rev *entity.ProblemRevision) error {
	problemID := rev.ProblemID
	if problemID == 0 {
		return entity.ErrCustomWrapper("problem id must be set for revision", nil)
	}

	query := fmt.Sprintf(`
		UPDATE %s SET outdated = true WHERE problemid = $1`, revisionTableName)
	// mark old revisions as outdated
	_, err := tx.Exec(query, problemID)
	if err != nil {
		return handleError(err)
	}

	query = fmt.Sprintf(`
		INSERT INTO %s
					( authorid,  problemid,  revision,  title,  memorylimit,  timelimit,  statementhtml,  statementpdf,  packagearchive,  outdated)
		VALUES(:authorid, :problemid, :revision, :title, :memorylimit, :timelimit, :statementhtml, :statementpdf, :packagearchive, false)
		RETURNING %s`, revisionTableName, revisionFieldsQuery)

	rows, err := tx.NamedQuery(query, rev)
	if err != nil {
		return handleError(err)
	}
	if rows.Next() {
		rows.StructScan(rev)
	}
	rows.Close()
	return nil
}

func (s ProblemRepository) Create(ctx context.Context, p *entity.Problem) error {
	if len(p.Revisions) != 1 {
		return fmt.Errorf("new problem should have 1 revision")
	}

	return TransactTxx(s.db, ctx, func(tx *sqlx.Tx) error {
		query := fmt.Sprintf(`
INSERT INTO %s ( authorid,  name,  externalurl)
VALUES         (:authorid, :name, :externalurl)
RETURNING %s`, problemTableName, problemFieldsQuery)

		rows, err := tx.NamedQuery(query, p)
		if err != nil {
			return handleError(err)
		}
		if rows.Next() {
			rows.StructScan(p)
		}
		rows.Close()

		for i := range p.Revisions {
			rev := &(p.Revisions[i])

			rev.ProblemID = p.ID

			// revision modification is not allowed. Revisions are immutable
			if rev.ID == 0 {
				err = s.createRevision(tx, rev)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// StatmentHtmlByRevisionId gets problem statment html for particular revision
func (s ProblemRepository) StatmentHtmlByProblemId(ctx context.Context, problemID entity.ID) (string, error) {
	var statementHtml string

	query := fmt.Sprintf(`SELECT rev.statementhtml FROM %s p LEFT JOIN %s rev ON rev.problemid = p.id AND rev.outdated = FALSE WHERE p.id = $1`, problemTableName, revisionTableName)

	// TODO remove
	log.Printf("%s \n vars: %v", query, problemID)

	err := s.db.QueryRowxContext(ctx, query, problemID).Scan(&statementHtml)
	if err != nil {
		return "", handleError(err)
	}

	return statementHtml, nil
}

// StatmentHtmlByRevisionId gets problem statment html for particular revision
func (s ProblemRepository) StatmentHtmlByRevisionId(ctx context.Context, revisionID entity.ID) (string, error) {
	var statementHtml string

	query := fmt.Sprintf(`SELECT statementhtml FROM %s WHERE id = $1`, revisionTableName)

	// TODO remove
	log.Printf("%s \n vars: %v", query, revisionID)

	err := s.db.QueryRowxContext(ctx, query, revisionID).Scan(&statementHtml)
	if err != nil {
		return "", handleError(err)
	}

	return statementHtml, nil
}

// StatmentPdfByRevisionId gets problem statment pdf for particular revision
func (s ProblemRepository) StatmentPdfByRevisionId(ctx context.Context, revisionID entity.ID) ([]byte, error) {
	var statementPdf []byte

	query := fmt.Sprintf(`SELECT statementpdf FROM %s WHERE id = $1`, revisionTableName)

	// TODO remove
	log.Printf("%s \n vars: %v", query, revisionID)

	err := s.db.QueryRowxContext(ctx, query, revisionID).Scan(&statementPdf)
	if err != nil {
		return nil, handleError(err)
	}

	return statementPdf, nil
}

// PackageArchiveByRevisionId gets problem statment pdf for particular revision
func (s ProblemRepository) PackageArchiveByRevisionId(ctx context.Context, revisionID entity.ID) ([]byte, error) {
	var packageArchive []byte

	query := fmt.Sprintf(`SELECT packagearchive FROM %s WHERE id = $1`, revisionTableName)

	err := s.db.QueryRowxContext(ctx, query, revisionID).Scan(&packageArchive)
	if err != nil {
		return nil, handleError(err)
	}

	return packageArchive, nil
}

func (s ProblemRepository) CanUserAccess(ctx context.Context, userID, problemID entity.ID) (contestIDs []entity.ID, err error) {
	err = s.db.SelectContext(ctx, &contestIDs, `
	SELECT cp.contestid
	FROM contest_problem cp
	INNER JOIN contest_user cu ON cp.contestid = cu.contestid
	INNER JOIN contest c ON cp.contestid = c.id
	WHERE cu.userid = $1 AND cp.problemid = $2 AND c.startTime < NOW();`, userID, problemID)
	if err != nil {
		return nil, handleError(err)
	}
	return
}

func (s ProblemRepository) AddRevision(ctx context.Context, newRevision *entity.ProblemRevision) error {
	return TransactTxx(s.db, ctx, func(tx *sqlx.Tx) error {
		return s.createRevision(tx, newRevision)
	})
}
