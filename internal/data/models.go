package data

import "database/sql"

type Models struct {
	Links LinkModel
}

func NewModelsInMemory(baseUrl string, suffixLength int) *Models {
	return &Models{
		Links: NewLinkModelInMemory(baseUrl, suffixLength),
	}
}

func NewModelsPostgres(baseUrl string, suffixLength int, db *sql.DB) *Models {
	return &Models{
		Links: &LinkModelPostgres{
			baseUrl:      baseUrl,
			suffixLength: suffixLength,
			DB:           db,
		},
	}
}
