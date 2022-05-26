package postgresql

import (
	"contestive/entity"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ContestRepository Store

const contestTableName = `"contest"`

var contestFieldNames = []string{"id", "createdat", "updatedat", "authorid", "contestname", "starttime", "endtime"}
var contestFieldsQuery = strings.Join(contestFieldNames, ", ")

func (s ContestRepository) Create(ctx context.Context, contest *entity.Contest) error {
	rows, err := s.db.NamedQueryContext(ctx, `
	INSERT INTO `+contestTableName+`
				( authorid,  contestname,  starttime,  endtime)
	VALUES(:authorid, :contestname, :starttime, :endtime)
	RETURNING `+contestFieldsQuery,
		contest)
	if err != nil {
		return handleError(err)
	}
	defer rows.Close()
	rows.Next()

	return handleError(rows.StructScan(contest))
}

func (s ContestRepository) GetByID(ctx context.Context, id entity.ID, userID entity.ID) (entity.Contest, error) {
	var contest entity.Contest
	var err error

	if userID != 0 {
		err = s.db.GetContext(ctx, &contest, `
		SELECT c.id, c.createdat, c.updatedat, c.authorid, c.contestname, c.starttime, c.endtime
		FROM `+contestTableName+` c
		INNER JOIN contest_user cu ON  c.id = cu.contestid AND cu.userid = $1
		WHERE c.starttime < NOW() AND id = $2
		LIMIT 1;`,
			userID, id,
		)
	} else {
		err = s.db.GetContext(ctx, &contest, `
		SELECT `+contestFieldsQuery+`
		FROM `+contestTableName+`
		WHERE id = $1
		LIMIT 1;`,
			id,
		)
	}
	if err != nil {
		return entity.Contest{}, handleError(err)
	}

	err = s.db.SelectContext(ctx, &contest.Problems, `SELECT problemid FROM contest_problem	WHERE contestid = $1;`, contest.ID)
	if err != nil {
		return entity.Contest{}, handleError(err)
	}

	if userID == 0 {
		err = s.db.SelectContext(ctx, &contest.Users, `SELECT userid FROM contest_user	WHERE contestid = $1;`, contest.ID)
		if err != nil {
			return entity.Contest{}, handleError(err)
		}
	}

	return contest, handleError(err)
}

// ListAll gets contests acording to options from all active contest list
func (s ContestRepository) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.Contest, int, error) {
	var contests []struct {
		entity.Contest
		Total int
	}

	qo := SqlListOptions(options, contestFieldNames)

	query := s.db.Rebind(fmt.Sprintf(`SELECT %s, count(*) OVER() AS total FROM %s WHERE %s %s %s`,
		contestFieldsQuery,
		contestTableName,
		qo.filterSql(),
		qo.orderSql(),
		qo.rangeSql(),
	))

	log.Printf("%s \n vars: %v", query, qo.filterVars)

	err := s.db.SelectContext(ctx, &contests, query, qo.filterVars...)

	if err != nil {
		return nil, 0, handleError(err)
	}

	res := make([]entity.Contest, len(contests))
	for i, v := range contests {
		res[i] = v.Contest
	}
	total := 0
	if len(contests) > 0 {
		total = contests[0].Total
	}

	return res, total, nil
}

// ListAll gets contests acording to options from all active contest list
func (s ContestRepository) ListAllForUsers(ctx context.Context, options entity.ListOptions, userid entity.ID) ([]entity.Contest, int, error) {
	var contests []struct {
		entity.Contest
		Total int
	}

	qo := SqlListOptions(options, contestFieldNames)

	query := s.db.Rebind(fmt.Sprintf(
		`SELECT c.id, c.createdat, c.updatedat, c.authorid, c.contestname, c.starttime, c.endtime, count(*) OVER() AS total
		FROM %s c
		INNER JOIN contest_user cu ON  c.id = cu.contestid AND cu.userid = %d
		WHERE c.starttime < NOW() AND c.endtime > NOW() AND %s %s %s`,
		contestTableName,
		userid,
		qo.filterSql(),
		qo.orderSql(),
		qo.rangeSql(),
	))

	log.Printf("%s \n vars: %v", query, qo.filterVars)

	err := s.db.SelectContext(ctx, &contests, query, qo.filterVars...)

	if err != nil {
		return nil, 0, handleError(err)
	}

	res := make([]entity.Contest, len(contests))
	for i, v := range contests {
		res[i] = v.Contest
	}
	total := 0
	if len(contests) > 0 {
		total = contests[0].Total
	}

	return res, total, nil
}

func (s ContestRepository) Update(ctx context.Context, contest *entity.Contest) error {
	return TransactTxx(s.db, ctx, func(tx *sqlx.Tx) error {

		err := s.updateContest(tx, contest)
		if err != nil {
			return err
		}
		err = s.updateUsers(tx, contest.ID, &contest.Users)
		if err != nil {
			return err
		}
		err = s.updateProblems(tx, contest.ID, &contest.Problems)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s ContestRepository) updateContest(tx *sqlx.Tx, contest *entity.Contest) error {
	rows, err := tx.NamedQuery(`
	UPDATE `+contestTableName+`
	SET ( contestname,  starttime,  endtime) =
      (:contestname, :starttime, :endtime)
	WHERE id = :id
	RETURNING `+contestFieldsQuery,
		contest)
	if err != nil {
		return handleError(err)
	}
	defer rows.Close()
	rows.Next()

	return handleError(rows.StructScan(contest))
}

func (s ContestRepository) updateUsers(tx *sqlx.Tx, contestID entity.ID, users *[]entity.ID) error {
	_, err := tx.Exec(`DELETE FROM contest_user	WHERE contestid = $1;`, contestID)
	if err != nil {
		return handleError(err)
	}

	if len(*users) == 0 {
		return nil
	}

	cid := strconv.Itoa(int(contestID))
	qb := strings.Builder{}

	qb.WriteString(`INSERT INTO contest_user ( contestid, userid ) VALUES `)
	for i, uid := range *users {
		if i != 0 {
			qb.WriteString(",(")
		} else {
			qb.WriteRune('(')
		}
		qb.WriteString(cid)
		qb.WriteRune(',')
		qb.WriteString(strconv.Itoa(int(uid)))
		qb.WriteRune(')')
	}
	qb.WriteRune(';')

	fmt.Println(qb.String())
	_, err = tx.Exec(qb.String())
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (s ContestRepository) updateProblems(tx *sqlx.Tx, contestID entity.ID, problems *[]entity.ID) error {
	_, err := tx.Exec(`DELETE FROM contest_problem	WHERE contestid = $1;`, contestID)
	if err != nil {
		return handleError(err)
	}

	if len(*problems) == 0 {
		return nil
	}

	cid := strconv.Itoa(int(contestID))
	qb := strings.Builder{}

	qb.WriteString(`INSERT INTO contest_problem ( contestid, problemid ) VALUES `)
	for i, pid := range *problems {
		if i != 0 {
			qb.WriteString(",(")
		} else {
			qb.WriteRune('(')
		}
		qb.WriteString(cid)
		qb.WriteRune(',')
		qb.WriteString(strconv.Itoa(int(pid)))
		qb.WriteRune(')')
	}
	qb.WriteRune(';')

	fmt.Println(qb.String())
	_, err = tx.Exec(qb.String())
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (s ContestRepository) Delete(ctx context.Context, id entity.ID) error {
	result, err := s.db.ExecContext(ctx, `
	DELETE FROM `+contestTableName+`
	WHERE id = $1`,
		id)
	if err != nil {
		return handleError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return handleError(err)
	}

	if rowsAffected == 0 {
		return entity.ErrNotFound(fmt.Errorf("contest delete afected 0 lines"))
	}

	return nil
}
