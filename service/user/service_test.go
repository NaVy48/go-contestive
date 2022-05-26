package user

import (
	"contestive/entity"
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
)

func newFixtureUser() entity.User {
	return entity.User{
		Entity: entity.Entity{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     "JS",
		PasswordHash: "123456789",
		FirstName:    "Ozzy",
		LastName:     "Osbourne",
	}
}

type passwordHasherStub struct{}

func (passwordHasherStub) HashPassword(p string) (string, error) {
	return p, nil
}

func Test_GetUserByUsername(t *testing.T) {
	is := is.New(t)
	repo := newInmem()
	m := NewService(repo, passwordHasherStub{})
	u := newFixtureUser()
	u.ID = 123
	repo.m[u.ID] = u

	retrieved, err := m.GetByUsername(context.Background(), u.Username)
	is.NoErr(err)
	is.Equal(retrieved, u)
}

// func Test_Create(t *testing.T) {
// 	repo := newInmem()
// 	m := NewService(repo)
// 	u := newFixtureUser()
// 	_, err := m.CreateUser(u.Username, u.Email, u.PasswordHash, u.FirstName, u.LastName)
// 	assert.Nil(t, err)
// 	assert.False(t, u.CreatedAt.IsZero())
// 	assert.True(t, u.UpdatedAt.IsZero())
// }

// func Test_SearchAndFind(t *testing.T) {
// 	repo := newInmem()
// 	m := NewService(repo)
// 	u1 := newFixtureUser()
// 	u2 := newFixtureUser()
// 	u2.FirstName = "Lemmy"

// 	uID, _ := m.CreateUser(u1.Username, u1.Email, u1.PasswordHash, u1.FirstName, u1.LastName)
// 	uID, _ := m.CreateUser(u1.Username, u1.Email, u1.PasswordHash, u1.FirstName, u1.LastName)
// 	_, _ = m.CreateUser(u2.Email, u2.Password, u2.FirstName, u2.LastName)

// 	t.Run("search", func(t *testing.T) {
// 		c, err := m.SearchUsers("ozzy")
// 		assert.Nil(t, err)
// 		assert.Equal(t, 1, len(c))
// 		assert.Equal(t, "Osbourne", c[0].LastName)

// 		c, err = m.SearchUsers("dio")
// 		assert.Equal(t, entity.ErrNotFound, err)
// 		assert.Nil(t, c)
// 	})
// 	t.Run("list all", func(t *testing.T) {
// 		all, err := m.ListUsers()
// 		assert.Nil(t, err)
// 		assert.Equal(t, 2, len(all))
// 	})

// 	t.Run("get", func(t *testing.T) {
// 		saved, err := m.GetUser(uID)
// 		assert.Nil(t, err)
// 		assert.Equal(t, u1.FirstName, saved.FirstName)
// 	})
// }

// func Test_Update(t *testing.T) {
// 	repo := newInmem()
// 	m := NewService(repo)
// 	u := newFixtureUser()
// 	saved, err := m.CreateUser(u.Username, u.Email, u.PasswordHash, u.FirstName, u.LastName)
// 	assert.Nil(t, err)
// 	saved, _ = m.GetUser(saved.ID)
// 	saved.FirstName = "Dio"
// 	assert.Nil(t, m.UpdateUser(&saved))
// 	updated, err := m.GetUser(saved.ID)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "Dio", updated.FirstName)
// 	assert.False(t, updated.UpdatedAt.IsZero())
// }

// func TestDelete(t *testing.T) {
// 	repo := newInmem()
// 	m := NewService(repo)
// 	u1 := newFixtureUser()
// 	u2 := newFixtureUser()
// 	savedU2, _ := m.CreateUser(u2.Username, u2.Email, u2.PasswordHash, u2.FirstName, u2.LastName)
// 	u2ID := savedU2.ID

// 	err := m.DeleteUser(u1.ID)
// 	assert.Equal(t, entity.ErrNotFound, err)

// 	err = m.DeleteUser(u2ID)
// 	assert.Nil(t, err)
// 	_, err = m.GetUser(u2ID)
// 	assert.Equal(t, entity.ErrNotFound, err)
// }
