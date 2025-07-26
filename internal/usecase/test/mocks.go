package test

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	entity "github.com/thisausername99/pantry-butler/internal/domain"
)

// MockRecipeRepository is a mock of RecipeRepository interface
type MockRecipeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRecipeRepositoryMockRecorder
}

// MockRecipeRepositoryMockRecorder is the mock recorder for MockRecipeRepository
type MockRecipeRepositoryMockRecorder struct {
	mock *MockRecipeRepository
}

// NewMockRecipeRepository creates a new mock instance
func NewMockRecipeRepository(ctrl *gomock.Controller) *MockRecipeRepository {
	mock := &MockRecipeRepository{ctrl: ctrl}
	mock.recorder = &MockRecipeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRecipeRepository) EXPECT() *MockRecipeRepositoryMockRecorder {
	return m.recorder
}

// GetRecipes mocks base method
func (m *MockRecipeRepository) GetRecipes(ctx context.Context) ([]*entity.Recipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecipes", ctx)
	ret0, _ := ret[0].([]*entity.Recipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecipes indicates an expected call of GetRecipes
func (mr *MockRecipeRepositoryMockRecorder) GetRecipes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecipes", reflect.TypeOf((*MockRecipeRepository)(nil).GetRecipes), ctx)
}

// GetRecipesByCuisine mocks base method
func (m *MockRecipeRepository) GetRecipesByCuisine(ctx context.Context, cuisine string) ([]*entity.Recipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecipesByCuisine", ctx, cuisine)
	ret0, _ := ret[0].([]*entity.Recipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecipesByCuisine indicates an expected call of GetRecipesByCuisine
func (mr *MockRecipeRepositoryMockRecorder) GetRecipesByCuisine(ctx, cuisine interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecipesByCuisine", reflect.TypeOf((*MockRecipeRepository)(nil).GetRecipesByCuisine), ctx, cuisine)
}

// MockPantryEntryRepository is a mock of PantryEntryRepository interface
type MockPantryEntryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPantryEntryRepositoryMockRecorder
}

// MockPantryEntryRepositoryMockRecorder is the mock recorder for MockPantryEntryRepository
type MockPantryEntryRepositoryMockRecorder struct {
	mock *MockPantryEntryRepository
}

// NewMockPantryEntryRepository creates a new mock instance
func NewMockPantryEntryRepository(ctrl *gomock.Controller) *MockPantryEntryRepository {
	mock := &MockPantryEntryRepository{ctrl: ctrl}
	mock.recorder = &MockPantryEntryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPantryEntryRepository) EXPECT() *MockPantryEntryRepositoryMockRecorder {
	return m.recorder
}

// GetPantryEntries mocks base method
func (m *MockPantryEntryRepository) GetPantryEntries(ctx context.Context) ([]*entity.PantryEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPantryEntries", ctx)
	ret0, _ := ret[0].([]*entity.PantryEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPantryEntries indicates an expected call of GetPantryEntries
func (mr *MockPantryEntryRepositoryMockRecorder) GetPantryEntries(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPantryEntries", reflect.TypeOf((*MockPantryEntryRepository)(nil).GetPantryEntries), ctx)
}

// InsertPantryEntry mocks base method
func (m *MockPantryEntryRepository) InsertPantryEntry(ctx context.Context, entry *entity.PantryEntryInput) (*entity.PantryEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPantryEntry", ctx, entry)
	ret0, _ := ret[0].(*entity.PantryEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertPantryEntry indicates an expected call of InsertPantryEntry
func (mr *MockPantryEntryRepositoryMockRecorder) InsertPantryEntry(ctx, entry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPantryEntry", reflect.TypeOf((*MockPantryEntryRepository)(nil).InsertPantryEntry), ctx, entry)
}
