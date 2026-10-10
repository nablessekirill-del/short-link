package service

import (
	"strings"
	"testing"
)

func TestNewCode(t *testing.T) {
	for _, n := range []int{1, 7, 16} {
		code, err := NewCode(n)
		if err != nil {
			t.Fatalf("NewCode(%d): %v", n, err)
		}
		if len(code) != n {
			t.Fatalf("want len %d, got %d (%q)", n, len(code), code)
		}
		for _, r := range code {
			if !strings.ContainsRune(alphabet, r) {
				t.Fatalf("unexpected char %q in %q", r, code)
			}
		}
	}
}

func TestNewCodeUnique(t *testing.T) {
	seen := make(map[string]struct{}, 10000)
	for i := 0; i < 10000; i++ {
		code, err := NewCode(7)
		if err != nil {
			t.Fatal(err)
		}
		if _, dup := seen[code]; dup {
			t.Fatalf("duplicate code %q after %d iterations", code, i)
		}
		seen[code] = struct{}{}
	}
}

func TestIsReserved(t *testing.T) {
	if !IsReserved("api") {
		t.Fatal("api must be reserved")
	}
	if IsReserved("abc1234") {
		t.Fatal("abc1234 must not be reserved")
	}
}
