package grpcx

import "context"

// AdminService is a placeholder admin RPC surface for proxyd.
type AdminService struct{}

func (AdminService) Ping(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return "pong", nil
}
