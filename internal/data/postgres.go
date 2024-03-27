package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "url.shortener/internal/proto"
)

type LinkModelPostgres struct {
	baseUrl string
	DB      *sql.DB
}

func (m *LinkModelPostgres) Insert(original *pb.OriginalUrl) (*pb.ShortUrl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	originalUrl := original.GetOriginalUrl()

	query := `
		SELECT 
			short_url 
		FROM 
			links 
		WHERE 
			original_url = $1
	`

	var shortUrl string
	err := m.DB.QueryRowContext(ctx, query, originalUrl).Scan(&shortUrl)
	if err == sql.ErrNoRows {
		query = `
			INSERT INTO links (original_url, short_url) 
			VALUES ($1, $2) 
		`

		suffix, err := generateUrlSuffix()
		if err != nil {
			return nil, err
		}

		shortUrl = fmt.Sprintf("%s/%s", m.baseUrl, suffix)

		_, err = m.DB.ExecContext(ctx, query, originalUrl, shortUrl)
		if err != nil {
			return nil, err
		}

		return &pb.ShortUrl{ShortUrl: shortUrl}, nil
	}

	if err != nil {
		return nil, err
	}

	return &pb.ShortUrl{ShortUrl: shortUrl}, nil
}

func (m *LinkModelPostgres) Get(short *pb.ShortUrl) (*pb.OriginalUrl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	shortUrl := short.GetShortUrl()

	query := `
		SELECT
			original_url
		FROM
			links
		WHERE
			short_url = $1
	`

	var originalUrl string

	err := m.DB.QueryRowContext(ctx, query, shortUrl).Scan(&originalUrl)
	if err == sql.ErrNoRows {
		return nil, ErrLinkNotFound
	}

	if err != nil {
		return nil, err
	}

	return &pb.OriginalUrl{OriginalUrl: originalUrl}, nil
}
