package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"friend-service/internal/auth-service/domain/auth/entities"
	"strings"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Create(ctx context.Context, attrs entities.AuthCredentialsAttrs) (entities.AuthCredentials, error) {
	var authCreds entities.AuthCredentials
	query := `INSERT INTO auth_users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id, email, password_hash`
	err := r.db.QueryRowContext(ctx, query, attrs.UserID, attrs.Email, attrs.HashedPassword).Scan(&authCreds.UserID, &authCreds.Email, &authCreds.HashedPassword)
	if err != nil {
		return entities.AuthCredentials{}, err
	}
	return authCreds, nil
}

func (r *AuthRepository) Get(ctx context.Context, filter entities.AuthCredentialsFilter) (entities.AuthCredentials, error) {
	var authCreds entities.AuthCredentials

	clauses := []string{}
	args := []any{}
	argIdx := 1

	if filter.UserID.IsPresent() {
		clauses = append(clauses, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, filter.UserID.MustGet())
		argIdx++
	}
	if filter.Email.IsPresent() {
		clauses = append(clauses, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, filter.Email.MustGet())
		argIdx++
	}
	if len(clauses) == 0 {
		return entities.AuthCredentials{}, fmt.Errorf("at least one filter must be provided")
	}

	query := `SELECT id, email, password_hash FROM auth_users WHERE ` + strings.Join(clauses, " AND ") + ` LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&authCreds.UserID, &authCreds.Email, &authCreds.HashedPassword)
	if err != nil {
		return entities.AuthCredentials{}, err
	}
	return authCreds, nil
}

func (r *AuthRepository) Update(ctx context.Context, userID string, attrs entities.AuthCredentialsAttrs) (entities.AuthCredentials, error) {
	var authCreds entities.AuthCredentials
	query := `UPDATE auth_users SET email = $1, password_hash = $2 WHERE id = $3 RETURNING id, email, password_hash`
	err := r.db.QueryRowContext(ctx, query, attrs.Email, attrs.HashedPassword, userID).Scan(&authCreds.UserID, &authCreds.Email, &authCreds.HashedPassword)
	if err != nil {
		return entities.AuthCredentials{}, err
	}
	return authCreds, nil
}

func (r *AuthRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM auth_users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
