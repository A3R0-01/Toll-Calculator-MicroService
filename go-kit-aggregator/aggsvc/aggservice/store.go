package aggservice

import (
	"fmt"

	"tollCalculator.com/types"
)

type MemoryStore struct {
	data map[int]float64
}

func (m *MemoryStore) Insert(d types.Distance) error {
	fmt.Println("this is coming from the bsns logic")
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	fmt.Println("this is coming from the bsns logic")
	distance, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obu id %d", id)
	}

	return distance, nil
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}

}
