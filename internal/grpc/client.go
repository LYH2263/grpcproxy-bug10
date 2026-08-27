package grpcx

import "context"

// Dial opens a placeholder client handle for diagnostics.
func Dial(ctx context.Context, target string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return target, nil
}
