package shp2geojson

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiPolygon(t *testing.T) {
	data := "D:\\temp\\MultiPolygon.shp"
	jsonStr, err := ConvertWithEncoding(data, "utf-8")
	assert.Equal(t, err, nil)
	//json := string(jsonStr)
	err = ioutil.WriteFile("./MultiPolygon.geojson", jsonStr, 0644)
	assert.Equal(t, err, nil)
}

func TestGBK(t *testing.T) {
	data := "data/GBKTest.shp"
	jsonStr, err := ConvertWithEncoding(data, "GBK")
	assert.Equal(t, err, nil)
	//json := string(jsonStr)
	err = ioutil.WriteFile("./gbktest.geojson", jsonStr, 0644)
	assert.Equal(t, err, nil)
}

func TestPoint(t *testing.T) {
	data := "./data/point.shp"
	jsonStr, err := Convert(data)
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile("./fixture/point.geojson", jsonStr, 0644)
	assert.Equal(t, err, nil)
}

func TestLine(t *testing.T) {
	data := "./data/line.shp"
	jsonStr, err := Convert(data)
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile("./fixture/line.geojson", jsonStr, 0644)
	assert.Equal(t, err, nil)
}

func TestPolygon(t *testing.T) {
	data := "./data/polygon.shp"
	jsonStr, err := Convert(data)
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile("./fixture/polygon.geojson", jsonStr, 0644)
	assert.Equal(t, err, nil)
}
