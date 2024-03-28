package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	pb "url.shortener/internal/proto"
	"url.shortener/internal/server"
)

func (app *application) newServer() *server.UrlShortenerServer {
	return server.NewServer(app.logger, app.models)
}

func (app *application) serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.config.ip, app.config.port))
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			app.loggingInterceptor,
			app.withRateLimitInterceptor(),
		),
	)

	pb.RegisterUrlShortenerServer(grpcServer, app.newServer())

	done := make(chan bool)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		grpcServer.GracefulStop()

		done <- true
	}()

	go func() {
		app.logger.PrintInfo("starting server", map[string]string{
			"address": fmt.Sprintf("%s:%d", app.config.ip, app.config.port),
			"storage": app.config.storage.storage_type,
			"limiter": fmt.Sprintf("enabled=%t, rps=%.2f, burst=%d", app.config.limiter.enabled, app.config.limiter.rps, app.config.limiter.burst),
		})

		err = grpcServer.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			app.logger.PrintFatal(err, nil)
		}
	}()

	<-done

	app.logger.PrintInfo("server gracefully stopped", nil)

	return nil
}
