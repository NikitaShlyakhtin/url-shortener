package main

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"url.shortener/internal/server"
)

func (app *application) withRateLimitInterceptor() grpc.ServerOption {
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

		resp, err := handler(ctx, req)
		if err != nil {
			if errors.Is(err, server.ErrServerError) {
				app.logger.PrintError(err, nil)
			}

			return nil, err
		}

		return resp, nil
	}

	return grpc.UnaryInterceptor(rateLimitInterceptor)
}
