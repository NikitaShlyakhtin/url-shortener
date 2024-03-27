package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	pb "url.shortener/internal/proto"
	"url.shortener/internal/server"
)

func (app *application) serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.config.ip, app.config.port))
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUrlShortenerServer(grpcServer, server.NewServer(app.logger))

	app.logger.PrintInfo("Starting server", map[string]string{
		"address": fmt.Sprintf("%s:%d", app.config.ip, app.config.port),
	})

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}