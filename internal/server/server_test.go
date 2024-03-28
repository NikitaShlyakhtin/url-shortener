package server

import (
	"context"
	"errors"
	"log"
	"testing"

	"url.shortener/internal/data"
	"url.shortener/internal/jsonlog"
	pb "url.shortener/internal/proto"
)

func TestCreateShortUrl(t *testing.T) {
	s := &UrlShortenerServer{
		models: *data.NewModelsInMemory("test", 6),
		logger: jsonlog.New(log.Writer(), jsonlog.LevelInfo),
	}

	ctx := context.Background()

	originalUrl := &pb.OriginalUrl{
		OriginalUrl: "example.com",
	}

	shortUrl, err := s.CreateShortUrl(ctx, originalUrl)
	if err != nil {
		t.Errorf("CreateShortUrl returned an error: %v", err)
	}

	if shortUrl.GetShortUrl() == "" {
		t.Error("CreateShortUrl returned an empty short URL")
	}
}

func TestGetOriginalUrl(t *testing.T) {
	s := &UrlShortenerServer{
		models: *data.NewModelsInMemory("test", 6),
		logger: jsonlog.New(log.Writer(), jsonlog.LevelInfo),
	}

	ctx := context.Background()

	original := &pb.OriginalUrl{
		OriginalUrl: "example.com",
	}

	short, err := s.CreateShortUrl(ctx, original)
	if err != nil {
		t.Errorf("CreateShortUrl returned an error: %v", err)
	}

	got, err := s.GetOriginalUrl(ctx, short)
	if err != nil {
		t.Errorf("GetOriginalUrl returned an error: %v", err)
	}

	if got.GetOriginalUrl() != original.GetOriginalUrl() {
		t.Errorf("GetOriginalUrl returned incorrect original URL, got: %s, want: %s", got.GetOriginalUrl(), original.GetOriginalUrl())
	}

	invalidShortUrl := &pb.ShortUrl{
		ShortUrl: "invalid",
	}

	_, err = s.GetOriginalUrl(ctx, invalidShortUrl)
	if err == nil {
		t.Error("GetOriginalUrl did not return an error for invalid short URL")
	} else if err != ErrNotFound {
		t.Errorf("GetOriginalUrl returned incorrect error, got: %v, want: %v", err, ErrNotFound)
	}
}

func TestValidateOriginalUrl(t *testing.T) {
	testCases := []struct {
		name        string
		originalUrl string
		expectedErr error
	}{
		{
			name:        "Valid URL",
			originalUrl: "example.com",
			expectedErr: nil,
		},
		{
			name:        "Empty URL",
			originalUrl: "",
			expectedErr: invalidArgumentError(errors.New("original URL must be provided")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateOriginalUrl(tc.originalUrl)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("validateOriginalUrl returned incorrect error, got: %v, want: %v", err, tc.expectedErr)
			}
		})
	}
}

func TestValidateShortUrl(t *testing.T) {
	testCases := []struct {
		name        string
		shortUrl    string
		expectedErr error
	}{
		{
			name:        "Valid Short URL",
			shortUrl:    "example.com",
			expectedErr: nil,
		},
		{
			name:        "Empty Short URL",
			shortUrl:    "",
			expectedErr: invalidArgumentError(errors.New("short URL must be provided")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateShortUrl(tc.shortUrl)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("validateShortUrl returned incorrect error, got: %v, want: %v", err, tc.expectedErr)
			}
		})
	}
}
