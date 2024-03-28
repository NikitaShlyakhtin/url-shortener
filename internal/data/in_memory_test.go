package data

import (
	"testing"

	pb "url.shortener/internal/proto"
)

func TestLinkModelInMemory_Insert(t *testing.T) {
	m := NewLinkModelInMemory("test", 6)

	original := &pb.OriginalUrl{
		OriginalUrl: "example.com",
	}

	short, err := m.Insert(original)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if short.ShortUrl == "" {
		t.Error("expected non-empty short URL")
	}

	short2, err := m.Insert(original)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if short2.ShortUrl != short.ShortUrl {
		t.Error("expected the same short URL for the same original URL")
	}
}

func TestLinkModelInMemory_Get(t *testing.T) {
	m := NewLinkModelInMemory("test", 6)

	original := &pb.OriginalUrl{
		OriginalUrl: "example.com",
	}

	short, err := m.Insert(original)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	result, err := m.Get(short)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.OriginalUrl != original.OriginalUrl {
		t.Errorf("expected original URL %s, got %s", original.OriginalUrl, result.OriginalUrl)
	}

	nonExistentShort := &pb.ShortUrl{
		ShortUrl: "nonexistent",
	}

	_, err = m.Get(nonExistentShort)
	if err != ErrLinkNotFound {
		t.Errorf("expected ErrLinkNotFound, got %v", err)
	}
}
