package data

type Models struct {
	Links LinkModel
}

func NewModelsInMemory() *Models {
	return &Models{
		Links: NewLinkModelInMemory(),
	}
}
