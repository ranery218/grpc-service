package ports

import "context"

type TokenGenerator interface {
	GenerateToken(ctx context.Context, nBytes int) (string, error)
}
