package postgres

import (
	"github.com/go-pg/pg/v10"
	"github.com/thisausername99/recipes-api/graph/model"
)

type PantryEntriesRepo struct {
	DB *pg.DB
}

func (m *PantryEntriesRepo) GetPantryEntries() ([]*model.PantryEntry, error) {
	var entries []*model.PantryEntry
	err := m.DB.Model(&entries).Select()
	if err != nil {
		return nil, err
	}
	return entries, nil
}
