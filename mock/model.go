package mock

import (
	"errors"
	"strconv"

	"github.com/target/gophersaurus/model"
)

// Model represents a model.Model.
type Model struct {
	ID    int
	data  string
	owner model.Model
}

// NewModel returns a *gf.Model.
func NewModel() *Model {
	return &Model{ID: 4037200794235010051, data: "mock model"}
}

// New implements gf.Model.
func (m *Model) New() model.Model {
	return NewModel()
}

// PathID implements gf.Model.
func (m *Model) PathID() string {
	return "model_id"
}

// SetID implements model.Model.
func (m *Model) SetID(id string) error {
	i, err := strconv.Atoi(id)
	if err == nil {
		m.ID = i
	}
	return err
}

// FindAll implements gf.Model.
func (m *Model) FindAll() ([]model.Model, error) {
	return []model.Model{m.New(), m.New(), m.New()}, nil
}

// FindByID implements gf.Model.
func (m *Model) FindByID(id string) error {
	return nil
}

// FindAllByOwner implements gf.Model.
func (m *Model) FindAllByOwner(owner model.Model) ([]model.Model, error) {
	models := []model.Model{m.New(), m.New(), m.New()}
	for _, model := range models {
		model.BelongsTo(owner)
	}
	return models, nil
}

// BelongsTo implements gf.Model.
func (m *Model) BelongsTo(owner model.Model) error {
	m.owner = owner
	return nil
}

// Save implements gf.Model.
func (m *Model) Save() error {
	return nil
}

// Delete implements gf.Model.
func (m *Model) Delete() error {
	return nil
}

// Validate implements gf.Model.
func (m *Model) Validate() error {
	return nil
}

// BaseModel represents a gf.Model.
type BaseModel struct {
	ID   int
	data string
}

// NewBaseModel returns a *gf.Model.
func NewBaseModel() *BaseModel {
	return &BaseModel{ID: 4037200794235010051, data: "mock model"}
}

// New implements gf.Model.
func (m *BaseModel) New() model.Model {
	return NewBaseModel()
}

// PathID implements gf.Model.
func (m *BaseModel) PathID() string {
	return "base_id"
}

// SetID implements gf.Model.
func (m *BaseModel) SetID(id string) error {
	return nil
}

// FindAll implements gf.Model.
func (m *BaseModel) FindAll() ([]model.Model, error) {
	return []model.Model{m.New(), m.New(), m.New()}, nil
}

// FindByID implements gf.Model.
func (m *BaseModel) FindByID(id string) error {
	i, err := strconv.Atoi(id)
	if err == nil {
		m.ID = i
	}
	return err
}

// FindAllByOwner implements gf.Model.
func (m *BaseModel) FindAllByOwner(owner model.Model) ([]model.Model, error) {
	return nil, errors.New("Base models usually do not have an owner.")
}

// BelongsTo implements gf.Model.
func (m *BaseModel) BelongsTo(owner model.Model) error {
	return errors.New("Base models usually do not have an owner.")
}

// Save implements gf.Model.
func (m *BaseModel) Save() error {
	return nil
}

// Delete implements gf.Model.
func (m *BaseModel) Delete() error {
	return nil
}

// Validate implements gf.Model.
func (m *BaseModel) Validate() error {
	return nil
}
