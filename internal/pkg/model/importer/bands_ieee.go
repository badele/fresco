package importer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"fresco/internal/pkg/model"

	"fresco/internal/pkg/tools"
)

const ieee_bands string = "http://www.classic.grss-ieee.org/fars3/barebones/eca_allocations_v1.js"

type IEEEBand struct {
	Region       string
	Khzmin       float64
	Khzmax       float64
	Primaryuse   string
	Secondaryuse string
	Footnotes    string
}

type IEEEBands []IEEEBand

// Get ressources list from Github repository
func (ieee *IEEEBands) GetRessouresList() []string {

	return []string{ieee_bands}
}

func (ieee *IEEEBands) ImportRessources(ressources []string, bands model.Bands) model.Bands {
	for _, url := range ressources {
		resp, err := tools.GetHttpContent(url)
		tools.CheckError(err)

		// Remove begin content
		var re = regexp.MustCompile(`(?s).*?\[ {`)
		var replace = "[ {"
		var count = 1
		resp = re.ReplaceAllStringFunc(resp, func(s string) string {
			if count == 0 {
				return s
			}

			count -= 1
			return re.ReplaceAllString(s, replace)
		})

		// Remove end content
		re = regexp.MustCompile(`(?s)];.*`)
		replace = "]"
		count = 1
		resp = re.ReplaceAllStringFunc(resp, func(s string) string {
			if count == 0 {
				return s
			}

			count -= 1
			return re.ReplaceAllString(s, replace)
		})

		// Replace region
		re = regexp.MustCompile(`(?s)"region": "11"`)
		replace = `"region": "Region1"`
		resp = re.ReplaceAllString(resp, replace)

		// Remove reference notes
		re = regexp.MustCompile(`(?s) \([0-9]+[^\)]+\)`)
		replace = ""
		resp = re.ReplaceAllString(resp, replace)

		// Remove (R) or (OR)
		re = regexp.MustCompile(`(?s) \(O?R\)`)
		replace = ""
		resp = re.ReplaceAllString(resp, replace)

		// Uppercase Primaryuse
		re = regexp.MustCompile(`"primaryuse": "(.*?)"`)
		resp = re.ReplaceAllStringFunc(resp, func(s string) string {
			slices := strings.Split(s, ":")
			if len(slices) != 2 {
				return s
			}
			return fmt.Sprintf(`"primaryuse": %s`, strings.ToUpper(slices[1]))
		})
		// Uppercase Secondaryuse
		re = regexp.MustCompile(`"secondaryuse": "(.*?)"`)
		resp = re.ReplaceAllStringFunc(resp, func(s string) string {
			slices := strings.Split(s, ":")
			if len(slices) != 2 {
				return s
			}
			return fmt.Sprintf(`"secondaryuse": %s`, strings.ToUpper(slices[1]))
		})

		// Remove footnote
		re = regexp.MustCompile(`(?s)"footnotes": ".*?"`)
		replace = `"footnotes": ""`
		resp = re.ReplaceAllString(resp, replace)

		// Convert to float
		re = regexp.MustCompile(`(?s)"khzmin": "([0-9]+.[0-9]+)"`)
		replace = `"khzmin": $1`
		resp = re.ReplaceAllString(resp, replace)
		re = regexp.MustCompile(`(?s)"khzmax": "([0-9]+.[0-9]+)"`)
		replace = `"khzmax": $1`
		resp = re.ReplaceAllString(resp, replace)

		// Unmarshal JSON content
		var ieee_bands []IEEEBand
		err = json.Unmarshal([]byte(resp), &ieee_bands)
		tools.CheckError(err)

		for _, band := range ieee_bands {
			salltypes := strings.ToLower(band.Primaryuse + "," + band.Secondaryuse)
			alltypes := strings.Split(salltypes, ",")
			for _, stype := range alltypes {
				stype = strings.Trim(stype, " ")

				switch {
				case strings.Contains(stype, "broadcasting"):
					stype = "broadcast"
				case strings.Contains(stype, "radiolocation"), strings.Contains(stype, "radionavigation"):
					stype = "radiolocation"
				case strings.Contains(stype, "meteorological aids"):
					stype = "meteorological"
				case strings.Contains(stype, "radio astronomy"):
					stype = "astronomy"
				case strings.Contains(stype, "standard frequency and time signal"):
					stype = "time"
				case strings.Contains(stype, "space"), strings.Contains(stype, "earth"):
					stype = "space"
				case strings.Contains(stype, "maritime"):
					stype = "marine"
				case strings.Contains(stype, "aeronautical"):
					stype = "aviation"
				case strings.Contains(stype, "mobile"), strings.Contains(stype, "fixed"), strings.Contains(stype, "not allocated"):
					stype = ""
				}

				newband := model.Band{
					Name:           "",
					Location:       "Region1",
					LowerFrequency: int64(band.Khzmin * 1000),
					UpperFrequency: int64(band.Khzmax * 1000),
					Source:         "IEEE",
					Author:         "",
					Type:           stype,
				}
				bands = append(bands, newband)
			}
		}
	}

	return bands
}
