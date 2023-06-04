package importer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"fresco/internal/pkg/model"

	"fresco/internal/pkg/tools"

	"github.com/gocolly/colly"
)

const bandsURL = "https://github.com/AlexandreRouma/SDRPlusPlus/tree/master/root/res/bandplans"
const rawMainURL = "https://raw.githubusercontent.com/AlexandreRouma/SDRPlusPlus/master/root/res/bandplans"

type SDRPPBands struct {
	Name        string `json:"name"`
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
	AuthorName  string `json:"author_name"`
	AuthorURL   string `json:"author_url"`
	Bands       []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Start int64  `json:"start"`
		End   int64  `json:"end"`
	} `json:"bands"`
}

// Get ressources list from Github repository
func (bands *SDRPPBands) GetRessouresList() []string {
	links := []string{}

	c := colly.NewCollector()

	r, _ := regexp.Compile(`.*/(.*?\.json)$`)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		if e.Attr("class") == "js-navigation-open Link--primary" {
			// e.Request.Visit(e.Attr("href"))
			link := e.Attr("href")

			if match := r.MatchString(link); match == true {
				filename := r.FindStringSubmatch(link)[1]
				absoluteURL := fmt.Sprintf("%s/%s", rawMainURL, filename)
				links = append(links, absoluteURL)
			}
		}
	})

	mainURL := fmt.Sprintf(bandsURL)
	c.Visit(mainURL)

	return links
}

// import ressources list from Github repository
func (sdrpp *SDRPPBands) ImportRessources(ressources []string, bands model.Bands) model.Bands {
	for _, url := range ressources {
		resp, err := tools.GetHttpContent(url)
		tools.CheckError(err)

		err = json.Unmarshal([]byte(resp), sdrpp)
		tools.CheckError(err)

		for _, line := range sdrpp.Bands {
			sname := ""
			switch {
			case strings.Contains(line.Name, "Radionav"):
			default:
				sname = line.Name
			}

			stype := line.Type
			switch {
			case strings.Contains(line.Type, "mobile"):
				stype = "phone"
			}

			newband := model.Band{
				Name:           sname,
				Location:       sdrpp.CountryName,
				LowerFrequency: line.Start,
				UpperFrequency: line.End,
				Source:         "SDR++",
				Author:         sdrpp.AuthorName,
				Type:           stype,
			}
			bands = append(bands, newband)
		}
	}

	return bands
}
