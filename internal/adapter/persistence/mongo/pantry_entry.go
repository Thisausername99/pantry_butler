package mongo

import (
	"github.com/go-pg/pg/v10"
	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

type PantryEntryRepo struct {
	DB *pg.DB
}

func (m *PantryEntryRepo) GetPantryEntries() ([]*entity.PantryEntry, error) {
	var entries []*entity.PantryEntry
	err := m.DB.Model(&entries).Select()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (m *PantryEntryRepo) InsertPantryEntry(entry *entity.PantryEntry) (*entity.PantryEntry, error) {
	_, err := m.DB.Model(entry).Returning("*").Insert()
	return entry, err
}
