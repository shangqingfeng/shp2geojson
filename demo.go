package main

import (
	"fmt"

	"github.com/shangqingfeng/shp2geojson"
)

func main() {
	json, _ := shp2geojson.ConvertWithEncoding("data/GBKTest.shp", "gbk")
	jsonStr := string(json)
	fmt.Println(jsonStr)
}
