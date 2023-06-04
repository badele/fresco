package main

import (
	"sort"

	"fresco/internal/pkg/model"

	"fresco/internal/pkg/model/importer"
)

func ImportBands() {
	newbands := model.Bands{
		model.Band{
			Name:           "",
			Location:       "",
			LowerFrequency: 0,
			UpperFrequency: 0,
			Source:         "",
			Type:           "",
		},
		model.Band{Name: "SLF", Location: "Worldwide", LowerFrequency: 3e+1, UpperFrequency: 3e+2, Source: "IEEE", Type: "band"},
		model.Band{Name: "ULF", Location: "Worldwide", LowerFrequency: 3e+2, UpperFrequency: 3e+3, Source: "IEEE", Type: "band"},
		model.Band{Name: "VLF", Location: "Worldwide", LowerFrequency: 3e+3, UpperFrequency: 3e+4, Source: "IEEE", Type: "band"},
		model.Band{Name: "LF", Location: "Worldwide", LowerFrequency: 3e+4, UpperFrequency: 3e+5, Source: "IEEE", Type: "band"},
		model.Band{Name: "MF", Location: "Worldwide", LowerFrequency: 3e+5, UpperFrequency: 3e+6, Source: "IEEE", Type: "band"},
		model.Band{Name: "HF", Location: "Worldwide", LowerFrequency: 3e+6, UpperFrequency: 3e+7, Source: "IEEE", Type: "band"},
		model.Band{Name: "VHF", Location: "Worldwide", LowerFrequency: 3e+7, UpperFrequency: 3e+8, Source: "IEEE", Type: "band"},
		model.Band{Name: "UHF", Location: "Worldwide", LowerFrequency: 3e+8, UpperFrequency: 3e+9, Source: "IEEE", Type: "band"},
		model.Band{Name: "SHF", Location: "Worldwide", LowerFrequency: 3e+9, UpperFrequency: 3e+10, Source: "IEEE", Type: "band"},
		model.Band{Name: "EHF", Location: "Worldwide", LowerFrequency: 3e+10, UpperFrequency: 3e+11, Source: "IEEE", Type: "band"},
		model.Band{Name: "THF", Location: "Worldwide", LowerFrequency: 3e+11, UpperFrequency: 3e+12, Source: "IEEE", Type: "band"},
	}

	var dbimporter importer.Importer

	// IEEE
	dbimporter = &importer.IEEEBands{}
	newbands = importer.Import(newbands, dbimporter)

	// SDRPP
	dbimporter = &importer.SDRPPBands{}
	newbands = importer.Import(newbands, dbimporter)

	// GQRX
	dbimporter = &importer.GQRXBands{}
	newbands = importer.Import(newbands, dbimporter)

	// Save
	sort.Sort(model.Bands(newbands))
	newbands.Write("datasets/bands.fwf")
}

func main() {
	ImportBands()
}
