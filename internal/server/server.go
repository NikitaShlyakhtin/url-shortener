package server

import (
	"context"
	"errors"

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
	if err := validateOriginalUrl(in.GetOriginalUrl()); err != nil {
		return nil, s.invalidArgumentError(err)
	}

	short, err := s.models.Links.Insert(in)
	if err != nil {
		return nil, s.serverError(err)
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
	if err := validateShortUrl(in.GetShortUrl()); err != nil {
		return nil, s.invalidArgumentError(err)
	}

	original, err := s.models.Links.Get(in)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrLinkNotFound):
			return nil, s.notFoundError()
		default:
			return nil, s.serverError(err)
		}
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

func validateOriginalUrl(originalUrl string) error {
	if originalUrl == "" {
		return errors.New("original URL must be provided")
	}

	return nil
}

func validateShortUrl(shortUrl string) error {
	if shortUrl == "" {
		return errors.New("short URL must be provided")
	}

	return nil
}
