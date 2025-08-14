package tokengen

import "context"

type TokenService interface {
	GenerateAccessToken(ctx context.Context, userID, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID, role string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, string, error) // Returns userID and role
}
