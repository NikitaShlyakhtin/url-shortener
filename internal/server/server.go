package server

import (
	"context"

	"url.shortener/internal/data"
	"url.shortener/internal/jsonlog"
	pb "url.shortener/internal/proto"
)

type UrlShortenerServer struct {
	logger *jsonlog.Logger
	models data.Models

	pb.UnimplementedUrlShortenerServer
}

func NewServer(logger *jsonlog.Logger, models data.Models) *UrlShortenerServer {
	return &UrlShortenerServer{
		logger: logger,
		models: models,
	}
}

func (s *UrlShortenerServer) CreateShortUrl(ctx context.Context, in *pb.OriginalUrl) (*pb.ShortUrl, error) {
	short, err := s.models.Links.Insert(in)
	if err != nil {
		return nil, err
	}

	shortUrl := short.GetShortUrl()

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
	original, err := s.models.Links.Get(in)
	if err != nil {
		return nil, err
	}

	originalUrl := original.GetOriginalUrl()

	r := &pb.OriginalUrl{
		OriginalUrl: originalUrl,
	}

	s.logger.PrintInfo("Original URL retrieved", map[string]string{
		"short_url":    in.GetShortUrl(),
		"original_url": originalUrl,
	})

	return r, nil
}
