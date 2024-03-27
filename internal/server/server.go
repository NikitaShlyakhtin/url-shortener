package server

import (
	"context"
	"fmt"

	"url.shortener/internal/jsonlog"
	pb "url.shortener/internal/proto"
)

type UrlShortenerServer struct {
	logger *jsonlog.Logger

	pb.UnimplementedUrlShortenerServer
}

func NewServer(logger *jsonlog.Logger) *UrlShortenerServer {
	return &UrlShortenerServer{logger: logger}
}

func (s *UrlShortenerServer) CreateShortUrl(ctx context.Context, in *pb.OriginalUrl) (*pb.ShortUrl, error) {
	shortUrl := fmt.Sprint("short.url/", len(in.GetOriginalUrl()))

	r := &pb.ShortUrl{
		ShortUrl: shortUrl,
	}

	s.logger.PrintInfo("Short URL created", map[string]string{
		"original_url": in.GetOriginalUrl(),
		"short_url":    shortUrl,
	})

	return r, nil
}

func (s *UrlShortenerServer) GetOriginalUrl(ctx context.Context, in *pb.ShortUrl) (*pb.OriginalUrl, error) {
	originalUrl := "example.com"

	r := &pb.OriginalUrl{
		OriginalUrl: originalUrl,
	}

	s.logger.PrintInfo("Original URL retrieved", map[string]string{
		"short_url":    in.GetShortUrl(),
		"original_url": originalUrl,
	})

	return r, nil
}
