package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

const (
	codeLength      = 7
	maxCodeAttempts = 5
	maxURLLength    = 2048
)

var (
	ErrInvalidURL     = errors.New("invalid url")
	ErrCodeGeneration = errors.New("could not generate unique")
)

type Service struct {
	storage LinkStorage
}

func New(storage LinkStorage) *Service {
	return &Service{storage: storage}
}

func validateURL(raw string) error {
	if raw == "" || len(raw) > maxURLLength {
		return ErrInvalidURL
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return ErrInvalidURL
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURL
	}
	if u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}

func (s *Service) Create(ctx context.Context, rawURL string) (*Link, error) {
	if err := validateURL(rawURL); err != nil {
		return nil, err
	}

	for i := 0; i < maxCodeAttempts; i++ {
		code, err := NewCode(codeLength)
		if err != nil {
			return nil, fmt.Errorf("generate code: %w", err)
		}
		if IsReserved(code) {
			continue
		}

		link := &Link{Code: code, URL: rawURL}
		err = s.storage.Create(ctx, link)
		if errors.Is(err, ErrCodeTaken) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("store link: %w", err)
		}
		return link, nil
	}

	return nil, ErrCodeGeneration
}

func (s *Service) Resolve(ctx context.Context, code string) (*Link, error) {
	link, err := s.storage.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, ErrNotFound
	}
	return link, nil
}

