package user

import (
	"errors"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")

// Repository defines storage operations for users.
type Repository interface {
	Create(user User) (User, error)
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	Update(id int, user User) (User, error)
	Delete(id int) error
}

// InMemoryRepository stores users in memory with thread safety.
type InMemoryRepository struct {
	mu     sync.Mutex
	items  map[int]User
	nextID int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items:  make(map[int]User),
		nextID: 1,
	}
}

func (r *InMemoryRepository) Create(user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.items[user.ID] = user
	return user, nil
}

func (r *InMemoryRepository) GetAll() ([]User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]User, 0, len(r.items))
	for _, user := range r.items {
		users = append(users, user)
	}

	return users, nil
}

func (r *InMemoryRepository) GetByID(id int) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.items[id]
	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryRepository) Update(id int, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return User{}, ErrUserNotFound
	}

	user.ID = id
	r.items[id] = user
	return user, nil
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrUserNotFound
	}

	delete(r.items, id)
	return nil
}
