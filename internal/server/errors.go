package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound    = status.Error(codes.NotFound, "the requested resource was not found")
	ErrRateLimit   = status.Error(codes.ResourceExhausted, "the rate limit for this resource has been exceeded")
	ErrPeerContext = status.Error(codes.Internal, "couldn't get peer from context")
)

func internalServerError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func invalidArgumentError(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}
