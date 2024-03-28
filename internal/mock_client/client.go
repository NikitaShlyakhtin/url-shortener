package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"url.shortener/internal/jsonlog"
	pb "url.shortener/internal/proto"
)

func main() {
	addr := "localhost:50051"

	logger := jsonlog.New(log.Writer(), jsonlog.LevelInfo)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer conn.Close()

	client := pb.NewUrlShortenerClient(conn)

	logger.PrintInfo("Client started", map[string]string{
		"address": addr,
	})

	for {
		var input string = "create;test.com"
		fmt.Scan(&input)
		split := strings.Split(input, ";")

		command := split[0]
		value := split[1]

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		switch command {
		case "create":
			shortUrl, err := client.CreateShortUrl(ctx, &pb.OriginalUrl{OriginalUrl: value})
			if err != nil {
				logger.PrintError(err, nil)
				continue
			}

			logger.PrintInfo("Short URL created", map[string]string{
				"original_url": value,
				"short_url":    shortUrl.GetShortUrl(),
			})
		case "get":
			originalUrl, err := client.GetOriginalUrl(ctx, &pb.ShortUrl{ShortUrl: value})
			if err != nil {
				logger.PrintError(err, nil)
				continue
			}

			logger.PrintInfo("Original URL retrieved", map[string]string{
				"short_url":    value,
				"original_url": originalUrl.GetOriginalUrl(),
			})
		}

		cancel()
	}

}
