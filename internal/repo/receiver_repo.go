package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type ReceiverRepository struct {
	db *sql.DB
}

func NewReceiverRepository(db *sql.DB) *ReceiverRepository {
	return &ReceiverRepository{db: db}
}

func (r *ReceiverRepository) Create(ctx context.Context, receiver *domain.Receiver) error {
	query := `
		INSERT INTO receivers (
			id, project_id, name, type, position_x, position_y, position_z,
			weight, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		receiver.ID, receiver.ProjectID, receiver.Name, receiver.Type,
		receiver.PositionX, receiver.PositionY, receiver.PositionZ,
		receiver.Weight, receiver.CreatedAt, receiver.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create receiver: %w", err)
	}
	
	return nil
}

func (r *ReceiverRepository) GetByID(ctx context.Context, id string) (*domain.Receiver, error) {
	query := `
		SELECT id, project_id, name, type, position_x, position_y, position_z,
			   weight, created_at, updated_at
		FROM receivers 
		WHERE id = ?
	`
	
	var receiver domain.Receiver
	var receiverType string
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&receiver.ID, &receiver.ProjectID, &receiver.Name, &receiverType,
		&receiver.PositionX, &receiver.PositionY, &receiver.PositionZ,
		&receiver.Weight, &receiver.CreatedAt, &receiver.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receiver not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get receiver: %w", err)
	}
	
	receiver.Type = domain.ReceiverType(receiverType)
	
	return &receiver, nil
}

func (r *ReceiverRepository) ListByProjectID(ctx context.Context, projectID string) ([]*domain.Receiver, error) {
	query := `
		SELECT id, project_id, name, type, position_x, position_y, position_z,
			   weight, created_at, updated_at
		FROM receivers 
		WHERE project_id = ?
		ORDER BY created_at
	`
	
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list receivers: %w", err)
	}
	defer rows.Close()
	
	var receivers []*domain.Receiver
	
	for rows.Next() {
		var receiver domain.Receiver
		var receiverType string
		
		err := rows.Scan(
			&receiver.ID, &receiver.ProjectID, &receiver.Name, &receiverType,
			&receiver.PositionX, &receiver.PositionY, &receiver.PositionZ,
			&receiver.Weight, &receiver.CreatedAt, &receiver.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan receiver: %w", err)
		}
		
		receiver.Type = domain.ReceiverType(receiverType)
		
		receivers = append(receivers, &receiver)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating receivers: %w", err)
	}
	
	return receivers, nil
}

func (r *ReceiverRepository) Update(ctx context.Context, receiver *domain.Receiver) error {
	query := `
		UPDATE receivers 
		SET name = ?, type = ?, position_x = ?, position_y = ?, position_z = ?,
		    weight = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		receiver.Name, receiver.Type, receiver.PositionX, receiver.PositionY, receiver.PositionZ,
		receiver.Weight, receiver.UpdatedAt, receiver.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update receiver: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("receiver not found: %s", receiver.ID)
	}
	
	return nil
}

func (r *ReceiverRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM receivers WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete receiver: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("receiver not found: %s", id)
	}
	
	return nil
}

func (r *ReceiverRepository) DeleteByProjectID(ctx context.Context, projectID string) error {
	query := `DELETE FROM receivers WHERE project_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete receivers by project: %w", err)
	}
	
	return nil
}
