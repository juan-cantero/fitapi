package services

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/juan-cantero/fitapi/internal/models"
	"github.com/juan-cantero/fitapi/internal/repositories"
)

func TestCreateEquipment(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		CreateFunc: func(ctx context.Context, eq *models.Equipment) error {
			// Simulate successful creation
			eq.ID = "test-id-123"
			return nil
		},
	}

	service := NewEquipmentService(mockRepo)

	req := &models.CreateEquipmentRequest{
		Name:        "Barbell",
		Description: "Olympic barbell",
	}

	equipment, err := service.CreateEquipment(context.Background(), "user-123", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if equipment.Name != "Barbell" {
		t.Errorf("Expected name 'Barbell', got '%s'", equipment.Name)
	}

	if equipment.UserID != "user-123" {
		t.Errorf("Expected userID 'user-123', got '%s'", equipment.UserID)
	}
}

func TestCreateEquipment_RepositoryError(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		CreateFunc: func(ctx context.Context, eq *models.Equipment) error {
			return errors.New("database error")
		},
	}

	service := NewEquipmentService(mockRepo)

	req := &models.CreateEquipmentRequest{
		Name: "Barbell",
	}

	_, err := service.CreateEquipment(context.Background(), "user-123", req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetEquipment_Success(t *testing.T) {
	expectedEquipment := &models.Equipment{
		ID:     "eq-1",
		Name:   "Dumbbell",
		UserID: "user-123",
	}

	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return expectedEquipment, nil
		},
	}

	service := NewEquipmentService(mockRepo)

	equipment, err := service.GetEquipment(context.Background(), "eq-1", "user-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if equipment.ID != "eq-1" {
		t.Errorf("Expected ID 'eq-1', got '%s'", equipment.ID)
	}
}

func TestGetEquipment_NotFound(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return nil, pgx.ErrNoRows
		},
	}

	service := NewEquipmentService(mockRepo)

	_, err := service.GetEquipment(context.Background(), "nonexistent", "user-123")

	if !errors.Is(err, ErrEquipmentNotFound) {
		t.Errorf("Expected ErrEquipmentNotFound, got %v", err)
	}
}

func TestGetEquipment_Unauthorized(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return &models.Equipment{
				ID:     "eq-1",
				UserID: "different-user",
			}, nil
		},
	}

	service := NewEquipmentService(mockRepo)

	_, err := service.GetEquipment(context.Background(), "eq-1", "user-123")

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized, got %v", err)
	}
}

func TestListEquipment(t *testing.T) {
	expectedList := []*models.Equipment{
		{ID: "eq-1", Name: "Barbell", UserID: "user-123"},
		{ID: "eq-2", Name: "Dumbbell", UserID: "user-123"},
	}

	mockRepo := &repositories.MockEquipmentRepository{
		FindAllFunc: func(ctx context.Context, userID string) ([]*models.Equipment, error) {
			if userID != "user-123" {
				return []*models.Equipment{}, nil
			}
			return expectedList, nil
		},
	}

	service := NewEquipmentService(mockRepo)

	list, err := service.ListEquipment(context.Background(), "user-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 2 {
		t.Errorf("Expected 2 items, got %d", len(list))
	}
}

func TestUpdateEquipment_Success(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return &models.Equipment{
				ID:     "eq-1",
				Name:   "Old Name",
				UserID: "user-123",
			}, nil
		},
		UpdateFunc: func(ctx context.Context, eq *models.Equipment) error {
			return nil
		},
	}

	service := NewEquipmentService(mockRepo)

	req := &models.UpdateEquipmentRequest{
		Name:        "New Name",
		Description: "Updated description",
	}

	updated, err := service.UpdateEquipment(context.Background(), "eq-1", "user-123", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.Name != "New Name" {
		t.Errorf("Expected name 'New Name', got '%s'", updated.Name)
	}
}

func TestUpdateEquipment_Unauthorized(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return &models.Equipment{
				ID:     "eq-1",
				UserID: "different-user",
			}, nil
		},
	}

	service := NewEquipmentService(mockRepo)

	req := &models.UpdateEquipmentRequest{Name: "New Name"}

	_, err := service.UpdateEquipment(context.Background(), "eq-1", "user-123", req)

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized, got %v", err)
	}
}

func TestDeleteEquipment_Success(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return &models.Equipment{
				ID:     "eq-1",
				UserID: "user-123",
			}, nil
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	service := NewEquipmentService(mockRepo)

	err := service.DeleteEquipment(context.Background(), "eq-1", "user-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteEquipment_Unauthorized(t *testing.T) {
	mockRepo := &repositories.MockEquipmentRepository{
		FindByIDFunc: func(ctx context.Context, id string) (*models.Equipment, error) {
			return &models.Equipment{
				ID:     "eq-1",
				UserID: "different-user",
			}, nil
		},
	}

	service := NewEquipmentService(mockRepo)

	err := service.DeleteEquipment(context.Background(), "eq-1", "user-123")

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized, got %v", err)
	}
}
