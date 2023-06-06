package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"fresco/internal/pkg/model"
	"fresco/internal/pkg/model/importer"
	"fresco/internal/pkg/ui"

	"github.com/alecthomas/kong"
)

func ImportBands() {
	oldbands := model.Bands{}
	oldbands.Read("datasets/bands.fwf")

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

	// Merge only news band data
	oldbands.InsertNewValue(&newbands)

	// Save
	sort.Slice(oldbands, func(i, j int) bool {
		return oldbands[i].UUID() < oldbands[j].UUID()
	})

	// Write result file
	oldbands.Write("datasets/bands.fwf")
}

var cli struct {
	Import struct {
		Dataset string `enum:"bands,channels" required:"" short:"d" help:"Dataset name [ bands|channels ]"`
	} `cmd:"" help:"Import dataset"`
}

func main() {
	if len(os.Args) >= 2 {
		ctx := kong.Parse(&cli)
		switch ctx.Command() {
		case "import":
			dataset := strings.ToLower(cli.Import.Dataset)
			switch dataset {
			case "bands":
				fmt.Printf("Importing %s ...\n", dataset)
				ImportBands()
			}
		default:
			fmt.Print(ctx.Command())
		}
	} else {
		app := ui.NewApp()
		app.Start()
	}
}
