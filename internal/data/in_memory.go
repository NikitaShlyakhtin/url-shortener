package data

import (
	"fmt"
	"sync"

	pb "url.shortener/internal/proto"
)

type LinkModelInMemory struct {
	baseUrl         string
	suffixLength    int
	originalToShort map[string]string
	shortToOriginal map[string]string
	mu              sync.RWMutex
}

func NewLinkModelInMemory(baseUrl string, suffixLength int) *LinkModelInMemory {
	return &LinkModelInMemory{
		baseUrl:         baseUrl,
		suffixLength:    suffixLength,
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
	}
}

func (m *LinkModelInMemory) Insert(original *pb.OriginalUrl) (*pb.ShortUrl, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	originalUrl := original.GetOriginalUrl()

	short, present := m.originalToShort[originalUrl]
	if present {
		return &pb.ShortUrl{ShortUrl: short}, nil
	}

	suffix, err := generateUrlSuffix(m.suffixLength)
	if err != nil {
		return nil, err
	}

	short = fmt.Sprintf("%s/%s", m.baseUrl, suffix)
	m.originalToShort[originalUrl] = short
	m.shortToOriginal[short] = originalUrl

	return &pb.ShortUrl{ShortUrl: short}, nil
}

func (m *LinkModelInMemory) Get(short *pb.ShortUrl) (*pb.OriginalUrl, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	shortUrl := short.GetShortUrl()

	original, present := m.shortToOriginal[shortUrl]
	if !present {
		return nil, ErrLinkNotFound
	}

	return &pb.OriginalUrl{OriginalUrl: original}, nil
}
