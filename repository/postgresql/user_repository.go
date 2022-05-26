package postgresql

import (
	"contestive/entity"
	"context"
	"fmt"
	"log"
	"strings"
)

type UserRepository Store

const userTableName = `"user"`

var userFieldNames = []string{"id", "username", "firstname", "lastname", "passwordhash", "createdat", "updatedat", "admin"}
var userFieldsQuery = strings.Join(userFieldNames, ", ")

func (s UserRepository) Create(ctx context.Context, user *entity.User) error {
	rows, err := s.db.NamedQueryContext(ctx, `
	INSERT INTO `+userTableName+`(username, firstname, lastname, passwordhash, admin)
	VALUES(:username, :firstname, :lastname, :passwordhash, :admin)
	RETURNING `+userFieldsQuery,
		user)
	if err != nil {
		return handleError(err)
	}
	defer rows.Close()
	rows.Next()

	return handleError(rows.StructScan(user))
}

func (s UserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	var u entity.User

	err := s.db.GetContext(ctx, &u, `
						SELECT `+userFieldsQuery+`
						FROM `+userTableName+`
						WHERE LOWER(username) = LOWER($1)
						LIMIT 1;`,
		username,
	)

	return u, handleError(err)
}

func (s UserRepository) GetByID(ctx context.Context, id entity.ID) (entity.User, error) {
	var u entity.User

	err := s.db.GetContext(ctx, &u, `
						SELECT `+userFieldsQuery+`
						FROM `+userTableName+`
						WHERE id = $1
						LIMIT 1;`,
		id,
	)

	return u, handleError(err)
}

// ListAll gets users according to options from all active user list
func (s UserRepository) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.User, int, error) {
	var users []struct {
		entity.User
		Total int
	}

	qo := SqlListOptions(options, userFieldNames)

	query := s.db.Rebind(fmt.Sprintf(`SELECT %s, count(*) OVER() AS total FROM %s WHERE %s %s %s`,
		userFieldsQuery,
		userTableName,
		qo.filterSql(),
		qo.orderSql(),
		qo.rangeSql(),
	))

	log.Printf("%s \n vars: %v", query, qo.filterVars)

	err := s.db.SelectContext(ctx, &users, query, qo.filterVars...)

	if err != nil {
		return nil, 0, handleError(err)
	}

	res := make([]entity.User, len(users))
	for i, v := range users {
		res[i] = v.User
	}
	total := 0
	if len(users) > 0 {
		total = users[0].Total
	}

	return res, total, nil
}

func (s UserRepository) Update(ctx context.Context, user *entity.User) error {
	rows, err := s.db.NamedQueryContext(ctx, `
	UPDATE `+userTableName+`
	SET ( username,  firstname,  lastname,  passwordhash,  admin,  updatedat) =
	    (:username, :firstname, :lastname, :passwordhash, :admin, :updatedat)
	WHERE id = :id
	RETURNING `+userFieldsQuery,
		user)
	if err != nil {
		return handleError(err)
	}
	defer rows.Close()
	rows.Next()

	return handleError(rows.StructScan(user))
}

func (s UserRepository) Delete(ctx context.Context, id entity.ID) error {
	result, err := s.db.ExecContext(ctx, `
	DELETE FROM `+userTableName+`
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
		return entity.ErrNotFound(fmt.Errorf("user delete afected 0 lines"))
	}

	return nil
}
