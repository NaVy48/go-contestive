package user

import (
	"context"
	"math/rand"

	"contestive/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]entity.User
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]entity.User{}
	return &inmem{
		m: m,
	}
}

// Create a user
func (r *inmem) Create(ctx context.Context, e *entity.User) error {
	e.ID = rand.Int63()
	r.m[e.ID] = *e
	return nil
}

// GetByID gets user by ID
func (r *inmem) GetByID(ctx context.Context, id entity.ID) (entity.User, error) {
	user, ok := r.m[id]
	if !ok {
		return user, entity.ErrNotFound(nil)
	}
	return user, nil
}

// GetByUsername gets user by ID
func (r *inmem) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	for _, v := range r.m {
		if v.Username == username {
			return v, nil
		}
	}
	return entity.User{}, entity.ErrNotFound(nil)
}

func (r *inmem) ListAll(ctx context.Context, options entity.ListOptions) ([]entity.User, int, error) {
	res := make([]entity.User, len(r.m))
	i := 0
	for _, v := range r.m {
		res[i] = v
		i++
	}
	return res, len(res), nil
}

// Update a user
func (r *inmem) Update(ctx context.Context, e *entity.User) error {
	_, err := r.GetByID(ctx, e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = *e
	return nil
}

//Delete a user
func (r *inmem) Delete(ctx context.Context, id entity.ID) error {
	if _, ok := r.m[id]; !ok {
		return entity.ErrNotFound(nil)
	}
	delete(r.m, id)
	return nil
}
