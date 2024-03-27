package main

import (
	"context"
	"log"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	originalUrl := "example.com"

	shortUrl, err := client.CreateShortUrl(ctx, &pb.OriginalUrl{OriginalUrl: originalUrl})
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("Short URL created", map[string]string{
		"original_url": originalUrl,
		"short_url":    shortUrl.GetShortUrl(),
	})

	_, err = client.GetOriginalUrl(ctx, &pb.ShortUrl{ShortUrl: shortUrl.GetShortUrl()})
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("Original URL retrieved", map[string]string{
		"short_url":    shortUrl.GetShortUrl(),
		"original_url": "example.com",
	})
}
