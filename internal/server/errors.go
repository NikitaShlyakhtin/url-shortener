package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UrlShortenerServer) invalidArgumentError(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}

func (s *UrlShortenerServer) serverError(err error) error {
	s.logger.PrintError(err, nil)

	message := "the server encountered a problem and could not process your request"

	return status.Error(codes.Internal, message)
}

func (s *UrlShortenerServer) notFoundError() error {
	return status.Error(codes.NotFound, "the requested resource was not found")
}
