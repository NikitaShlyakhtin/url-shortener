package data

import (
	"errors"
	"fmt"
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

type LinkModelInMemory struct {
	baseUrl         string
	originalToShort map[string]string
	shortToOriginal map[string]string
}

func NewLinkModelInMemory(baseUrl string) *LinkModelInMemory {
	return &LinkModelInMemory{
		baseUrl:         baseUrl,
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
	}
}

func (m *LinkModelInMemory) Insert(original *pb.OriginalUrl) (*pb.ShortUrl, error) {
	originalUrl := original.GetOriginalUrl()

	short, present := m.originalToShort[originalUrl]
	if present {
		return &pb.ShortUrl{ShortUrl: short}, nil
	}

	sufix, err := generateUrlSufix()
	if err != nil {
		return nil, err
	}

	short = fmt.Sprintf("%s/%s", m.baseUrl, sufix)
	m.originalToShort[originalUrl] = short
	m.shortToOriginal[short] = originalUrl

	return &pb.ShortUrl{ShortUrl: short}, nil
}

func (m *LinkModelInMemory) Get(short *pb.ShortUrl) (*pb.OriginalUrl, error) {
	shortUrl := short.GetShortUrl()

	original, present := m.shortToOriginal[shortUrl]
	if !present {
		return nil, ErrLinkNotFound
	}

	return &pb.OriginalUrl{OriginalUrl: original}, nil
}

func generateUrlSufix() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	const shortURLLength = 10

	shortURL := make([]byte, shortURLLength)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortURL), nil
}
