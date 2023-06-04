package model

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/o1egl/fwencoder"
)

/*
##########################################################
# Band
##########################################################
*/

type Band struct {
	Name           string
	Location       string
	LowerFrequency int64
	UpperFrequency int64
	Mode           string
	Step           string
	Source         string
	Author         string
	Type           string
}

func (band *Band) AndFiltered(filter string) bool {
	selected := true

	filters := strings.Split(filter, "&")
	for _, search := range filters {

		// if previous selected result is false stop loop
		if !selected {
			return false
		}

		r := regexp.MustCompile(`^([^\=]+)\=([A-Za-z0-9].*)`)
		if !r.MatchString(search) {
			selected = true
			continue
		}

		groups := r.FindStringSubmatch(search)
		fieldname := groups[1]
		searchword := groups[2]

		value := reflect.ValueOf(band).Elem()
		fieldcontent := value.FieldByName(fieldname)

		if fieldcontent.IsValid() {
			selected = selected && strings.Contains(strings.ToLower(fieldcontent.String()), strings.ToLower(searchword))
		}
	}

	return selected
}

func (band *Band) IsFiltered(filter string) bool {
	selected := false

	// Or filter
	filters := strings.Split(filter, "|")
	for _, search := range filters {
		selected = selected || band.AndFiltered(search)
	}

	return selected
}

/*
##########################################################
# Bands
##########################################################
*/

type Bands []Band

func (bands *Bands) Print() {
	b, err := fwencoder.Marshal(bands)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func (bands *Bands) Read(filename string) {
	if _, err := os.Stat(filename); err == nil {

		f, _ := os.Open(filename)
		defer f.Close()

		// var datas Bands
		err := fwencoder.UnmarshalReader(f, bands)

		if err != nil {
			panic(err)
		}
	}
}

func (bands *Bands) Write(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	_ = fwencoder.MarshalWriter(f, bands)
}

// Sorting
func (bands Bands) Len() int      { return len(bands) }
func (bands Bands) Swap(i, j int) { bands[i], bands[j] = bands[j], bands[i] }
func (bands Bands) Less(i, j int) bool {
	if bands[i].LowerFrequency != bands[j].LowerFrequency {
		return bands[i].LowerFrequency < bands[j].LowerFrequency
	}

	return bands[i].UpperFrequency < bands[j].UpperFrequency
}

func NewBands(filename string) Bands {
	bands := Bands{}
	bands.Read(filename)

	return bands
}
