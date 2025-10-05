package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juan-cantero/fitapi/internal/models"
)

// EquipmentRepository defines the interface for equipment data access
type EquipmentRepository interface {
	Create(ctx context.Context, equipment *models.Equipment) error
	FindByID(ctx context.Context, id string) (*models.Equipment, error)
	FindAll(ctx context.Context, userID string) ([]*models.Equipment, error)
	Update(ctx context.Context, equipment *models.Equipment) error
	Delete(ctx context.Context, id string) error
}

// PostgresEquipmentRepository is the PostgreSQL implementation of EquipmentRepository
type PostgresEquipmentRepository struct {
	db *pgxpool.Pool
}

// NewPostgresEquipmentRepository creates a new PostgreSQL equipment repository
func NewPostgresEquipmentRepository(db *pgxpool.Pool) EquipmentRepository {
	return &PostgresEquipmentRepository{db: db}
}

// Create inserts a new equipment record into the database
func (r *PostgresEquipmentRepository) Create(ctx context.Context, equipment *models.Equipment) error {
	equipment.ID = uuid.New().String()

	query := `
		INSERT INTO equipment (id, name, description, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		equipment.ID,
		equipment.Name,
		equipment.Description,
		equipment.UserID,
	).Scan(&equipment.CreatedAt, &equipment.UpdatedAt)

	return err
}

// FindByID retrieves a single equipment by ID
func (r *PostgresEquipmentRepository) FindByID(ctx context.Context, id string) (*models.Equipment, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM equipment
		WHERE id = $1
	`

	equipment := &models.Equipment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&equipment.ID,
		&equipment.Name,
		&equipment.Description,
		&equipment.UserID,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return equipment, nil
}

// FindAll retrieves all equipment for a specific user
func (r *PostgresEquipmentRepository) FindAll(ctx context.Context, userID string) ([]*models.Equipment, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM equipment
		WHERE user_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipmentList []*models.Equipment
	for rows.Next() {
		equipment := &models.Equipment{}
		err := rows.Scan(
			&equipment.ID,
			&equipment.Name,
			&equipment.Description,
			&equipment.UserID,
			&equipment.CreatedAt,
			&equipment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		equipmentList = append(equipmentList, equipment)
	}

	return equipmentList, rows.Err()
}

// Update updates an existing equipment record
func (r *PostgresEquipmentRepository) Update(ctx context.Context, equipment *models.Equipment) error {
	query := `
		UPDATE equipment
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		equipment.Name,
		equipment.Description,
		equipment.ID,
	).Scan(&equipment.UpdatedAt)

	return err
}

// Delete removes an equipment record from the database
func (r *PostgresEquipmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM equipment WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
