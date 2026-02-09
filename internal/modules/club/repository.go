package club

import (
	"errors"
	"sync"
)

var ErrClubNotFound = errors.New("club not found")

// Repository defines storage operations for clubs.
type Repository interface {
	Create(club Club) (Club, error)
	GetAll() ([]Club, error)
	GetByID(id int) (Club, error)
	Update(id int, club Club) (Club, error)
	Delete(id int) error
}

// InMemoryRepository stores clubs in memory with thread safety.
type InMemoryRepository struct {
	mu     sync.Mutex
	items  map[int]Club
	nextID int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items:  make(map[int]Club),
		nextID: 1,
	}
}

func (r *InMemoryRepository) Create(club Club) (Club, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	club.ID = r.nextID
	r.nextID++
	r.items[club.ID] = club
	return club, nil
}

func (r *InMemoryRepository) GetAll() ([]Club, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	clubs := make([]Club, 0, len(r.items))
	for _, club := range r.items {
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (r *InMemoryRepository) GetByID(id int) (Club, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	club, ok := r.items[id]
	if !ok {
		return Club{}, ErrClubNotFound
	}
	return club, nil
}

func (r *InMemoryRepository) Update(id int, club Club) (Club, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return Club{}, ErrClubNotFound
	}

	club.ID = id
	r.items[id] = club
	return club, nil
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrClubNotFound
	}

	delete(r.items, id)
	return nil
}
