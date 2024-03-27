package data

type Models struct {
	Links LinkModel
}

func NewModelsInMemory(baseUrl string) *Models {
	return &Models{
		Links: NewLinkModelInMemory(baseUrl),
	}
}
