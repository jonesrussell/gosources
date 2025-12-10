package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/jonesrussell/gosources/internal/models"
)

// GlobalSelectorsRepository handles database operations for global selectors
type GlobalSelectorsRepository struct {
	db *sqlx.DB
}

// NewGlobalSelectorsRepository creates a new GlobalSelectorsRepository
func NewGlobalSelectorsRepository(db *sqlx.DB) *GlobalSelectorsRepository {
	return &GlobalSelectorsRepository{db: db}
}

// Get retrieves the global selectors configuration
// Returns the first (and should be only) record
func (r *GlobalSelectorsRepository) Get(ctx context.Context) (*models.GlobalSelectors, error) {
	var gs models.GlobalSelectors
	query := `SELECT id, selectors, created_at, updated_at FROM global_selectors LIMIT 1`

	err := r.db.GetContext(ctx, &gs, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("global selectors not found")
		}
		return nil, err
	}

	return &gs, nil
}

// Update updates the global selectors configuration
func (r *GlobalSelectorsRepository) Update(ctx context.Context, selectors *models.SelectorConfig) error {
	query := `
		UPDATE global_selectors
		SET selectors = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = (SELECT id FROM global_selectors LIMIT 1)
		RETURNING id, created_at, updated_at
	`

	var gs models.GlobalSelectors
	err := r.db.GetContext(ctx, &gs, query, selectors)
	if err != nil {
		return err
	}

	return nil
}

// Create creates a new global selectors record (should only be used for initialization)
func (r *GlobalSelectorsRepository) Create(ctx context.Context, selectors *models.SelectorConfig) (*models.GlobalSelectors, error) {
	query := `
		INSERT INTO global_selectors (selectors)
		VALUES ($1)
		RETURNING id, selectors, created_at, updated_at
	`

	var gs models.GlobalSelectors
	err := r.db.GetContext(ctx, &gs, query, selectors)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}
