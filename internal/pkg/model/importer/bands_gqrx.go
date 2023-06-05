package importer

import (
	"encoding/csv"
	"fresco/internal/pkg/model"
	"fresco/internal/pkg/tools"
	"io"
	"strconv"
	"strings"
)

const gqrx_bands string = "https://raw.githubusercontent.com/gqrx-sdr/gqrx/master/resources/bandplan.csv"

type GQRXBand struct {
	MinFrequency float64
	MaxFrequency float64
	Mode         string
	Step         string
	Color        string
	Name         string
}

type GQRXBands []GQRXBand

// Get ressources list from Github repository
func (gqrx *GQRXBands) GetRessouresList() []string {

	return []string{gqrx_bands}
}

func (gqrx *GQRXBands) ImportRessources(ressources []string, bands model.Bands) model.Bands {
	for _, url := range ressources {
		resp, err := tools.GetHttpContent(url)
		tools.CheckError(err)

		reader := csv.NewReader(strings.NewReader(resp))
		reader.Comment = '#'

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			tools.CheckError(err)

			minfreq, _ := strconv.ParseInt(strings.Trim(record[0], " "), 10, 64)
			maxfreq, _ := strconv.ParseInt(strings.Trim(record[1], " "), 10, 64)

			stype := ""
			switch {
			case strings.Contains(record[5], "Ham Band"):
				stype = "amateur"
			case strings.Contains(record[5], "Broadcast"):
				stype = "broadcast"
			case strings.Contains(record[5], "Air Band"):
				stype = "aviation"
			case strings.Contains(record[5], "Marine"):
				stype = "marine"
			case strings.Contains(record[5], "Radionavigation"):
				stype = "radionavigation"
			case strings.Contains(record[5], "Military Air"):
				stype = "military"
			case strings.Contains(record[5], "Mil Sat"):
				stype = "military"
			default:
				stype = "utility"
			}

			newband := model.Band{
				Name:           strings.Trim(record[5], " "),
				Location:       "WorldWide",
				LowerFrequency: minfreq,
				UpperFrequency: maxfreq,
				Mode:           strings.Trim(record[2], " "),
				Source:         "GQRX",
				Author:         "",
				Type:           stype,
			}
			bands = append(bands, newband)
		}
	}
	return bands
}
