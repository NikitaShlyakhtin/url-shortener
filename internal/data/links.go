package data

import (
	"errors"
	"math/rand"

	pb "url.shortener/internal/proto"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type LinkModel interface {
	Insert(original *pb.OriginalUrl) (*pb.ShortUrl, error)
	Get(short *pb.ShortUrl) (*pb.OriginalUrl, error)
}

func generateUrlSuffix() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	const shortURLLength = 10

	shortURL := make([]byte, shortURLLength)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortURL), nil
}
