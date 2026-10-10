package service

import (
	"context"
	"errors"
)

var (
	ErrNotFound  = errors.New("link not found")
	ErrCodeTaken = errors.New("code already taken")
)

type LinkStorage interface {
	Create(ctx context.Context, link *Link) error
	GetByCode(ctx context.Context, code string) (*Link, error)
}
