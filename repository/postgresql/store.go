package postgresql

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sqlx.DB
}

// UserRepository return user store
func (s Store) UserRepository() UserRepository {
	return UserRepository(s)
}

// ProblemRepository return problem store
func (s Store) ProblemRepository() ProblemRepository {
	return ProblemRepository(s)
}

// ContestRepository return problem store
func (s Store) ContestRepository() ContestRepository {
	return ContestRepository(s)
}

// SubmissionRepository return problem store
func (s Store) SubmissionRepository() SubmissionRepository {
	return SubmissionRepository(s)
}

// // ContestRepository return user store
// func (s *Store) ContestRepository(ctx context.Context) repository.ContestRepository {
// 	return &contestRepository{s.db, ctx}
// }

// // SubmissionRepository return user store
// func (s *Store) SubmissionRepository(ctx context.Context) repository.SubmissionRepository {
// 	return &submissionRepository{s.db, ctx}
// }

// Connect connects to a database
func Connect(address, username, password, database string) (Store, error) {
	var s Store

	connstr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?connect_timeout=10",
		username, password, address, database,
	)

	db, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		return s, err
	}
	db.SetMaxOpenConns(20)

	s = Store{db}

	err = s.Migrate()
	if err != nil {
		return s, err
	}

	return s, nil
}

// Migrate migrates the store database.
func (s *Store) Migrate() error {
	for _, q := range migrate {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

// Migrate migrates the store database.
func (s *Store) Seed() error {
	for _, q := range seed {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

// Drop drops the store database.
func (s *Store) Drop() error {
	for _, q := range drop {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

// Reset resets the store database.
func (s *Store) Reset() error {
	err := s.Drop()
	if err != nil {
		return err
	}

	err = s.Migrate()
	if err != nil {
		return err
	}

	err = s.Seed()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
