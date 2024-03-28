package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"url.shortener/internal/server"
)

func (app *application) withRateLimitInterceptor() grpc.UnaryServerInterceptor {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	rateLimitInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		pr, ok := peer.FromContext(ctx)
		if !ok {
			return nil, server.ErrPeerContext
		}

		ip := strings.Split(pr.Addr.String(), ":")[0]

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.logger.PrintError(server.ErrRateLimit, map[string]string{"ip": ip})
			return nil, server.ErrRateLimit
		}

		mu.Unlock()

		return handler(ctx, req)
	}

	return rateLimitInterceptor
}

func (app *application) loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	app.logger.PrintInfo("gRPC request", map[string]string{
		"method":  info.FullMethod,
		"request": fmt.Sprintf("%+v", req),
	})

	resp, err := handler(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Internal {
			app.logger.PrintFatal(err, nil)
		}

		app.logger.PrintError(err, nil)
	} else {
		app.logger.PrintInfo("gRPC response", map[string]string{
			"method":   info.FullMethod,
			"response": fmt.Sprintf("%+v", resp),
		})
	}

	return resp, err
}
