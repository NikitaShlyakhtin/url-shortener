package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrServerError = status.Error(codes.Internal, "the server encountered a problem and could not process your request")
	ErrNotFound    = status.Error(codes.NotFound, "the requested resource was not found")
	ErrRateLimit   = status.Error(codes.ResourceExhausted, "the rate limit for this resource has been exceeded")
	ErrPeerContext = status.Error(codes.Internal, "couldn't get peer from context")
)

func (s *UrlShortenerServer) invalidArgumentError(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}

func (s *UrlShortenerServer) serverError(err error) error {
	s.logger.PrintError(err, nil)

	return ErrServerError
}

func (s *UrlShortenerServer) notFoundError() error {
	return ErrNotFound
}
