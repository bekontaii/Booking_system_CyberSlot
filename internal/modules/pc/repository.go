package pc

import (
	"errors"
	"sync"
)

var ErrPCNotFound = errors.New("pc not found")

type Repository interface {
	Create(pc PC) (PC, error)
	GetAll() ([]PC, error)
	GetByID(id int) (PC, error)
	GetByClubID(clubID int) ([]PC, error)
	Update(id int, pc PC) (PC, error)
	Delete(id int) error
}

type InMemoryRepository struct {
	mu     sync.Mutex
	items  map[int]PC
	nextID int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items:  make(map[int]PC),
		nextID: 1,
	}
}

func (r *InMemoryRepository) Create(pc PC) (PC, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pc.ID = r.nextID
	r.nextID++
	r.items[pc.ID] = pc
	return pc, nil
}

func (r *InMemoryRepository) GetAll() ([]PC, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pcs := make([]PC, 0, len(r.items))
	for _, pc := range r.items {
		pcs = append(pcs, pc)
	}
	return pcs, nil
}

func (r *InMemoryRepository) GetByID(id int) (PC, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pc, ok := r.items[id]
	if !ok {
		return PC{}, ErrPCNotFound
	}

	return pc, nil
}

func (r *InMemoryRepository) GetByClubID(clubID int) ([]PC, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pcs := make([]PC, 0)
	for _, pc := range r.items {
		if pc.ClubID == clubID {
			pcs = append(pcs, pc)
		}
	}

	return pcs, nil
}

func (r *InMemoryRepository) Update(id int, pc PC) (PC, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return PC{}, ErrPCNotFound
	}

	pc.ID = id
	r.items[id] = pc
	return pc, nil
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrPCNotFound
	}

	delete(r.items, id)
	return nil
}
