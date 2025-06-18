package ai

import "context"

type AIProvider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
