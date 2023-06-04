package importer

import (
	"fresco/internal/pkg/model"
)

type Importer interface {
	GetRessouresList() []string // Get a ressouces to Download
	ImportRessources(ressources []string, bands model.Bands) model.Bands
}

func Import(bands model.Bands, i Importer) model.Bands {
	res := i.GetRessouresList()
	bands = i.ImportRessources(res, bands)

	return bands
}
