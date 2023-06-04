package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToFrequency_result(t *testing.T) {
	value, err := ConvertStringToFrequency("2.437 GHz")
	assert.Equal(t, value, float64(2437000000))
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("7.325265 MHz")
	assert.Equal(t, value, float64(7325265))
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("123.456 KHz")
	assert.Equal(t, value, float64(123456))
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("2345 Hz")
	assert.Equal(t, value, float64(2345))
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("2345")
	assert.Equal(t, value, float64(2345))
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("1 Hz")
	assert.Equal(t, value, float64(1))
	assert.Nil(t, err)
}

func TestConvertStringToFrequency_error(t *testing.T) {
	value, err := ConvertStringToFrequency("ghz")

	assert.Equal(t, value, float64(0))
	assert.EqualError(t, err, "Can't decode Frequency value")
}

func TestConvertFrequencyToString(t *testing.T) {
	value := ConvertFrequencyToString(2437000000)
	assert.Equal(t, value, "2.437000000 GHz")

	value = ConvertFrequencyToString(7325265)
	assert.Equal(t, value, "7.325265 MHz")

	value = ConvertFrequencyToString(2345)
	assert.Equal(t, value, "2.345 KHz")

	value = ConvertFrequencyToString(0)
	assert.Equal(t, value, "0 Hz")
}
