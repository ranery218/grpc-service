package postrges

import (
	"context"
	"database/sql"
	"fmt"
	"friend-service/internal/user-service/domain/user/entities"
	"strings"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(ctx context.Context, attrs entities.ProfileAttrs) (entities.Profile, error) {
	var profile entities.Profile
	query := `INSERT INTO user_profiles (user_id, username) VALUES ($1, $2) RETURNING user_id, username`

	err := r.db.QueryRowContext(ctx, query, attrs.ID, attrs.Username).Scan(&profile.ID, &profile.Username)
	if err != nil {
		return entities.Profile{}, err
	}

	return profile, nil
}

func (r *ProfileRepository) Get(ctx context.Context, filter entities.ProfileFilter) (entities.Profile, error) {
	var profile entities.Profile

	query := `SELECT user_id, username FROM user_profiles WHERE `

	filterClauses := []string{}
	args := []any{}
	argIndex := 1

	if filter.ID.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filter.ID.MustGet())
		argIndex++
	}
	if filter.Username.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, filter.Username.MustGet())
	}

	if len(filterClauses) == 0 {
		return entities.Profile{}, entities.ErrEmptyFilter
	}

	query += strings.Join(filterClauses, " AND ")
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&profile.ID, &profile.Username)
	if err != nil {
		return entities.Profile{}, err
	}

	return profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, filter entities.ProfileFilter, attrs entities.ProfileAttrs) (entities.Profile, error) {
	var profile entities.Profile

	query := `UPDATE user_profiles SET username = $1 WHERE user_id = $2 RETURNING user_id, username`

	err := r.db.QueryRowContext(ctx, query, attrs.Username, filter.ID.MustGet()).Scan(&profile.ID, &profile.Username)
	if err != nil {
		return entities.Profile{}, err
	}

	return profile, nil
}

func (r *ProfileRepository) GetAll(ctx context.Context, filter entities.ProfileFilter) ([]entities.Profile, error) {
	var profiles []entities.Profile

	query := `SELECT user_id, username FROM user_profiles`

	filterClauses := []string{}
	args := []any{}
	argIndex := 1

	if filter.Username.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, filter.Username.MustGet())
	}

	query += " WHERE " + strings.Join(filterClauses, " AND ")

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var profile entities.Profile
		if err := rows.Scan(&profile.ID, &profile.Username); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *ProfileRepository) Delete(ctx context.Context, filter entities.ProfileFilter) error {
	query := `DELETE FROM user_profiles WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, filter.ID.MustGet())
	if err != nil {
		return err
	}

	return nil
}