package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/juan-cantero/fitapi/internal/models"
	"github.com/juan-cantero/fitapi/internal/repositories"
)

var (
	ErrEquipmentNotFound = errors.New("equipment not found")
	ErrUnauthorized      = errors.New("unauthorized to perform this action")
)

// EquipmentService handles business logic for equipment
type EquipmentService struct {
	repo repositories.EquipmentRepository
}

// NewEquipmentService creates a new equipment service
func NewEquipmentService(repo repositories.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

// CreateEquipment creates a new equipment for a user
func (s *EquipmentService) CreateEquipment(ctx context.Context, userID string, req *models.CreateEquipmentRequest) (*models.Equipment, error) {
	equipment := &models.Equipment{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
	}

	if err := s.repo.Create(ctx, equipment); err != nil {
		return nil, fmt.Errorf("failed to create equipment: %w", err)
	}

	return equipment, nil
}

// GetEquipment retrieves a single equipment by ID
func (s *EquipmentService) GetEquipment(ctx context.Context, id string, userID string) (*models.Equipment, error) {
	equipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEquipmentNotFound
		}
		return nil, fmt.Errorf("failed to get equipment: %w", err)
	}

	// Check ownership
	if equipment.UserID != userID {
		return nil, ErrUnauthorized
	}

	return equipment, nil
}

// ListEquipment retrieves all equipment for a user
func (s *EquipmentService) ListEquipment(ctx context.Context, userID string) ([]*models.Equipment, error) {
	equipment, err := s.repo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list equipment: %w", err)
	}

	return equipment, nil
}

// UpdateEquipment updates an existing equipment
func (s *EquipmentService) UpdateEquipment(ctx context.Context, id string, userID string, req *models.UpdateEquipmentRequest) (*models.Equipment, error) {
	// First check if equipment exists and user owns it
	equipment, err := s.GetEquipment(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	equipment.Name = req.Name
	equipment.Description = req.Description

	if err := s.repo.Update(ctx, equipment); err != nil {
		return nil, fmt.Errorf("failed to update equipment: %w", err)
	}

	return equipment, nil
}

// DeleteEquipment deletes an equipment
func (s *EquipmentService) DeleteEquipment(ctx context.Context, id string, userID string) error {
	// First check if equipment exists and user owns it
	if _, err := s.GetEquipment(ctx, id, userID); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete equipment: %w", err)
	}

	return nil
}
