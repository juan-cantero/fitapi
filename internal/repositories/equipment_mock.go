package repositories

import (
	"context"

	"github.com/juan-cantero/fitapi/internal/models"
)

// MockEquipmentRepository is a mock implementation for testing
type MockEquipmentRepository struct {
	CreateFunc  func(ctx context.Context, equipment *models.Equipment) error
	FindByIDFunc func(ctx context.Context, id string) (*models.Equipment, error)
	FindAllFunc  func(ctx context.Context, userID string) ([]*models.Equipment, error)
	UpdateFunc   func(ctx context.Context, equipment *models.Equipment) error
	DeleteFunc   func(ctx context.Context, id string) error
}

func (m *MockEquipmentRepository) Create(ctx context.Context, equipment *models.Equipment) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, equipment)
	}
	return nil
}

func (m *MockEquipmentRepository) FindByID(ctx context.Context, id string) (*models.Equipment, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockEquipmentRepository) FindAll(ctx context.Context, userID string) ([]*models.Equipment, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, userID)
	}
	return []*models.Equipment{}, nil
}

func (m *MockEquipmentRepository) Update(ctx context.Context, equipment *models.Equipment) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, equipment)
	}
	return nil
}

func (m *MockEquipmentRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
