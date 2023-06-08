package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToFrequency_result(t *testing.T) {
	value, err := ConvertStringToFrequency("2.437 GHz")
	assert.Equal(t, int64(2437000000), value)
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("7.325265 MHz")
	assert.Equal(t, int64(7325265), value)
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("123.456 KHz")
	assert.Equal(t, int64(123456), value)
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("2345 Hz")
	assert.Equal(t, int64(2345), value)
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("2345")
	assert.Equal(t, int64(2345), value)
	assert.Nil(t, err)

	value, err = ConvertStringToFrequency("1 Hz")
	assert.Equal(t, int64(1), value)
	assert.Nil(t, err)
}

func TestConvertStringToFrequency_error(t *testing.T) {
	value, err := ConvertStringToFrequency("ghz")

	assert.Equal(t, int64(0), value)
	assert.EqualError(t, err, "Can't decode Frequency value")
}

func TestConvertFrequencyToString(t *testing.T) {
	value := ConvertFrequencyToString(2437000000)
	assert.Equal(t, "2.437000000 GHz", value)

	value = ConvertFrequencyToString(7325265)
	assert.Equal(t, "7.325265 MHz", value)

	value = ConvertFrequencyToString(2345)
	assert.Equal(t, "2.345 KHz", value)

	value = ConvertFrequencyToString(0)
	assert.Equal(t, "0 Hz", value)
}
