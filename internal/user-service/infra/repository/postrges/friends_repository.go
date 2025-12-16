package postrges

import (
	"context"
	"database/sql"
	"friend-service/internal/user-service/domain/user/entities"

	"github.com/lib/pq"
)

type FriendsRepository struct {
	db *sql.DB
}

func NewFriendsRepository(db *sql.DB) *FriendsRepository {
	return &FriendsRepository{db: db}
}

func (r *FriendsRepository) Create(ctx context.Context, attrs entities.FriendsAttrs) (entities.Friends, error) {
	var friends entities.Friends

	query := `INSERT INTO user_friends (user_id, friend_ids) VALUES ($1, $2) RETURNING user_id, friend_ids`

	err := r.db.QueryRowContext(ctx, query, attrs.UserID, pq.Array(attrs.Friends)).Scan(&friends.UserID, pq.Array(&friends.Friends))
	if err != nil {
		return entities.Friends{}, err
	}

	return friends, nil
}

func (r *FriendsRepository) Get(ctx context.Context, filter entities.FriendsFilter) (entities.Friends, error) {
	var friends entities.Friends

	query := `SELECT user_id, friend_ids FROM user_friends WHERE user_id = $1`

	err := r.db.QueryRowContext(ctx, query, filter.UserID.MustGet()).Scan(&friends.UserID, pq.Array(&friends.Friends))
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Friends{}, entities.ErrUsersFriendsNotFound
		}
		return entities.Friends{}, err
	}

	return friends, nil
}

func (r *FriendsRepository) Update(ctx context.Context, filter entities.FriendsFilter, attrs entities.FriendsAttrs) (entities.Friends, error) {
	var friends entities.Friends

	query := `UPDATE user_friends SET friend_ids = $1 WHERE user_id = $2 RETURNING user_id, friend_ids`

	err := r.db.QueryRowContext(ctx, query, pq.Array(attrs.Friends), filter.UserID.MustGet()).Scan(&friends.UserID, pq.Array(&friends.Friends))
	if err != nil {
		return entities.Friends{}, err
	}

	return friends, nil
}
