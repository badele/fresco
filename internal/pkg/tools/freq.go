package tools

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Frequency float64

const (
	Hz  Frequency = 1.0
	KHz           = Hz * 1000
	MHz           = KHz * 1000
	GHz           = MHz * 1000
)

const FreqValidationRegex = `^([0-9]+\.?[0-9]*) ?(Hz|KHz|MHz|GHz)?$`

// Frequency is in Mgz
func ConvertFrequencyToString(freq int64) string {
	switch {
	case freq > 1000000000:
		return fmt.Sprintf("%.9f GHz", float64(freq)/float64(GHz))
	case freq > 1000000:
		return fmt.Sprintf("%.6f MHz", float64(freq)/float64(MHz))
	case freq > 1000:
		return fmt.Sprintf("%.3f KHz", float64(freq)/float64(KHz))
	}

	return fmt.Sprintf("%d Hz", freq)
}

// Frequency is in Mgz
func ConvertStringToFrequency(ifreq string) (int64, error) {
	r := regexp.MustCompile(FreqValidationRegex)
	if !r.MatchString(ifreq) {
		return 0, errors.New("Can't decode Frequency value")
	}

	groups := r.FindStringSubmatch(ifreq)

	freq, err := strconv.ParseFloat(groups[1], 64)
	CheckError(err)
	unit := groups[2]

	var multiplier Frequency = Hz
	switch {
	case unit == "KHz":
		multiplier = KHz
	case unit == "MHz":
		multiplier = MHz
	case unit == "GHz":
		multiplier = GHz
	}

	return int64(freq * float64(multiplier)), nil
}
