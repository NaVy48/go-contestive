package postgresql

import (
	"contestive/entity"
	"contestive/internal/api/params"
	"contestive/repository"
	"context"
	"database/sql"
	"strings"
)

type testRepository struct {
	db  *sql.DB
	ctx context.Context
}

const testTableName = `"test"`

var testFieldNames = []string{"id", "problem_id", "order_num", "title", "input", "expected_output", "created_at", "deleted_at"}
var testFieldsQuery = strings.Join(testFieldNames, ", ")

func (s *testRepository) scan(scanner scanner) (*repository.Test, error) {
	t := new(repository.Test)
	err := scanner.Scan(
		&t.TestID,
		&t.ProblemID,
		&t.OrderNum,
		&t.Title,
		&t.Input,
		&t.ExpectedOutput,
		&t.CreatedAt,
		&t.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}
	return t, nil
}

// Get all problems
func (s *testRepository) List() ([]*repository.Test, int, error) {
	q := params.QueryFromContext(s.ctx)
	tests := make([]*repository.Test, 0, 50)

	rows, err := s.db.QueryContext(s.ctx, `
		SELECT id, problem_id "problemId", order_num "orderNum", title, LEFT(input, 100), LEFT(expected_output, 100), created_at, deleted_at
		FROM `+testTableName+`
		WHERE deleted_at IS NULL `+q.Filter.SqlAllowedFieldsWithAnd(map[string]string{"id": "id", "problemId": "problem_id"})+`
		`+q.Sort.Sql()+q.Range.Sql()+`;`,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		problem, err := s.scan(rows)
		if err != nil {
			return nil, 0, err
		}
		tests = append(tests, problem)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// FIXME fix total
	return tests, 5, nil
}

func (s *testRepository) SingleByID(id int64) (*repository.Test, error) {
	return s.scan(s.db.QueryRowContext(s.ctx, `
			SELECT `+testFieldsQuery+`
			FROM `+testTableName+`
			WHERE id = $1 AND deleted_at IS NULL
			LIMIT 1`,
		id,
	))
}

func (s *testRepository) Create(t repository.Test) (*repository.Test, error) {
	return s.scan(s.db.QueryRowContext(s.ctx, `
	INSERT INTO `+testTableName+`(problem_id, order_num, title, input, expected_output, created_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING `+testFieldsQuery,
		t.ProblemID,
		t.OrderNum,
		t.Title,
		t.Input,
		t.ExpectedOutput,
		t.CreatedAt,
	))
}
func (s *testRepository) Update(id int64, t repository.Test) (*repository.Test, error) {
	return s.scan(s.db.QueryRowContext(s.ctx, `
	UPDATE `+testTableName+`
	SET
	order_num = $2, title = $3
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING `+testFieldsQuery,
		id,
		t.OrderNum,
		t.Title,
	))
}

func (s *testRepository) Delete(id int64) (*repository.Test, error) {
	return s.scan(s.db.QueryRowContext(s.ctx, `
	UPDATE `+testTableName+`
	SET
	deleted_at = NOW()
	WHERE id = $1
	RETURNING `+testFieldsQuery,
		id,
	))
}
