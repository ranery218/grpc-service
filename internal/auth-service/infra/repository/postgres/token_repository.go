package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"friend-service/internal/auth-service/domain/auth/entities"
	"strings"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, attrs entities.RefreshTokenAttrs) (entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, token_hash, expires_at`
	err := r.db.QueryRowContext(ctx, query, attrs.ID, attrs.UserID, attrs.HashedRefreshToken, attrs.ExpiresAt).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.HashedRefreshToken, &refreshToken.ExpiresAt)
	if err != nil {
		return entities.RefreshToken{}, err
	}
	return refreshToken, nil
}

func (r *TokenRepository) Get(ctx context.Context, filter entities.RefreshTokenFilter) (entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	query := `SELECT id, user_id, token_hash, expires_at FROM refresh_tokens WHERE `

	filterClauses := []string{}
	args := []any{}
	argIndex := 1

	if filter.ID.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, filter.ID.MustGet())
		argIndex++
	}
	if filter.UserID.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filter.UserID.MustGet())
		argIndex++
	}
	if filter.HashedRefreshToken.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("token_hash = $%d", argIndex))
		args = append(args, filter.HashedRefreshToken.MustGet())
	}

	query += fmt.Sprintf("%s LIMIT 1", strings.Join(filterClauses, " AND "))

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.HashedRefreshToken, &refreshToken.ExpiresAt)
	if err != nil {
		return entities.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *TokenRepository) Delete(ctx context.Context, filter entities.RefreshTokenFilter) error {
	query := `DELETE FROM refresh_tokens WHERE `

	filterClauses := []string{}
	args := []any{}
	argIndex := 1

	if filter.ID.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, filter.ID.MustGet())
		argIndex++
	}
	if filter.UserID.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filter.UserID.MustGet())
		argIndex++
	}
	if filter.HashedRefreshToken.IsPresent() {
		filterClauses = append(filterClauses, fmt.Sprintf("token_hash = $%d", argIndex))
		args = append(args, filter.HashedRefreshToken.MustGet())
	}

	query += strings.Join(filterClauses, " AND ")

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
