package service_test

import (
	"context"
	"errors"
	"shortener/internal/service"
	"shortener/internal/storage"
	"testing"
	"time"
)

func TestCreateAndResolve(t *testing.T) {
	ctx := context.Background()
	svc := service.New(storage.NewMemory())

	link, err := svc.Create(ctx, "https://example.com/page")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(link.Code) != 7 {
		t.Fatalf("want code len 7, got %q", link.Code)
	}

	got, err := svc.Resolve(ctx, link.Code)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if got.URL != "https://example.com/page" {
		t.Fatalf("unexpected url %q", got.URL)
	}
}

func TestCreateInvalidURL(t *testing.T) {
	svc := service.New(storage.NewMemory())

	cases := []string{
		"",
		"not a url",
		"ftp://example.com",
		"javascript:alert(1)",
		"https://",
	}
	for _, c := range cases {
		if _, err := svc.Create(context.Background(), c); !errors.Is(err, service.ErrInvalidURL) {
			t.Errorf("url %q: want ErrInvalidURL, got %v", c, err)
		}
	}
}

func TestResolveNotFound(t *testing.T) {
	svc := service.New(storage.NewMemory())
	if _, err := svc.Resolve(context.Background(), "nope"); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestResolveExpired(t *testing.T) {
	ctx := context.Background()
	st := storage.NewMemory()
	svc := service.New(st)

	past := time.Now().Add(-time.Hour)
	st.Create(ctx, &service.Link{Code: "old1234", URL: "https://example.com", ExpiresAt: &past})

	if _, err := svc.Resolve(ctx, "old1234"); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("want ErrNotFound for expired link, got %v", err)
	}
}

type flaky struct {
	service.LinkStorage
	failures int
}

func (f *flaky) Create(ctx context.Context, l *service.Link) error {
	if f.failures > 0 {
		f.failures--
		return service.ErrCodeTaken
	}
	return f.LinkStorage.Create(ctx, l)
}

func TestCreateRetriesOnCollision(t *testing.T) {
	svc := service.New(&flaky{LinkStorage: storage.NewMemory(), failures: 2})
	if _, err := svc.Create(context.Background(), "https://example.com"); err != nil {
		t.Fatalf("want succes after retries, got %v", err)
	}
}

func TestCreateGivesUpAftertooManyCollisions(t *testing.T) {
	svc := service.New(&flaky{LinkStorage: storage.NewMemory(), failures: 100})
	_, err := svc.Create(context.Background(), "https://example.com")
	if !errors.Is(err, service.ErrCodeGeneration) {
		t.Fatalf("want ErrCodeGeneration, got %v", err)
	}
}
