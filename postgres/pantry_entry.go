package postgres

import (
	"github.com/go-pg/pg/v10"
	"github.com/thisausername99/recipes-api/models"
)

type PantryEntryRepo struct {
	DB *pg.DB
}

func (m *PantryEntryRepo) GetPantryEntries() ([]*models.PantryEntry, error) {
	var entries []*models.PantryEntry
	err := m.DB.Model(&entries).Select()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (m *PantryEntryRepo) InsertPantryEntry(entry *models.PantryEntry) (*models.PantryEntry, error) {
	_, err := m.DB.Model(entry).Returning("*").Insert()
	return entry, err
}
