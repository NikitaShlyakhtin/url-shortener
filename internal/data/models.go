package data

import "database/sql"

type Models struct {
	Links LinkModel
}

func NewModelsInMemory(baseUrl string) *Models {
	return &Models{
		Links: NewLinkModelInMemory(baseUrl),
	}
}

func NewModelsPostgres(baseUrl string, db *sql.DB) *Models {
	return &Models{
		Links: &LinkModelPostgres{
			baseUrl: baseUrl,
			DB:      db,
		},
	}
}
