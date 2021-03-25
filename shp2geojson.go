package shp2geojson

import (
	"fmt"
	"reflect"

	geojson "github.com/paulmach/go.geojson"
	"github.com/shangqingfeng/go-shp"
)

func Convert(input string) ([]byte, error) {
	return ConvertWithEncoding(input, "utf-8")
}

// Convert shapefile to geojson
// currently not support Multi Geometry
func ConvertWithEncoding(input string, encoding string) ([]byte, error) {
	shape, err := shp.OpenWithEncoding(input, encoding)
	if err != nil {
		return nil, err
	}
	defer shape.Close()

	// fields from the attribute table (DBF)
	fields := shape.Fields()
	fc := geojson.NewFeatureCollection()
	// loop through all features in the shapefile
	for shape.Next() {
		n, s := shape.Shape()
		var feature *geojson.Feature

		switch s.(type) {
		case *shp.Point:
			pt, _ := s.(*shp.Point)
			feature = geojson.NewPointFeature([]float64{pt.X, pt.Y})
		case *shp.PolyLine:
			line, _ := s.(*shp.PolyLine)
			// if line.NumParts != 1 {
			// 	fmt.Println("Warning: more than 1 part in polyline!")
			// }
			coordinates := make([][][]float64, line.NumParts)
			var i int32
			for i = 0; i < line.NumParts; i++ {

				var startIndex, endIndex int32
				startIndex = line.Parts[i]
				if i == line.NumParts-1 {
					endIndex = int32(len(line.Points))
				} else {
					endIndex = line.Parts[i+1]
				}
				coordinates[i] = make([][]float64, endIndex-startIndex)
				cnt := 0
				for j := startIndex; j < endIndex; j = j + 1 {
					coordinates[i][cnt] = make([]float64, 2)
					coordinates[i][cnt][0] = line.Points[j].X
					coordinates[i][cnt][1] = line.Points[j].Y
					cnt = cnt + 1
				}
			}
			feature = geojson.NewMultiLineStringFeature(coordinates...)
		case *shp.Polygon:
			polygon, _ := s.(*shp.Polygon)

			coordinates := make([][][]float64, polygon.NumParts)
			var i int32
			for i = 0; i < polygon.NumParts; i = i + 1 {
				var startIndex, endIndex int32
				startIndex = polygon.Parts[i]
				if i == polygon.NumParts-1 {
					endIndex = int32(len(polygon.Points))
				} else {
					endIndex = polygon.Parts[i+1]
				}

				coordinates[i] = make([][]float64, endIndex-startIndex)
				cnt := 0
				for j := startIndex; j < endIndex; j = j + 1 {
					coordinates[i][cnt] = make([]float64, 2)
					coordinates[i][cnt][0] = polygon.Points[j].X
					coordinates[i][cnt][1] = polygon.Points[j].Y
					cnt = cnt + 1
				}
			}
			var lastAllCoordinates = make([][][][]float64, 1)
			lastAllCoordinates[0] = make([][][]float64, 1)
			lastAllCoordinates[0][0] = make([][]float64, 1)
			for i := 0; i < len(coordinates[0]); i++ {
				lastAllCoordinates[0][0] = coordinates[0]
			}
			if len(coordinates) > 1 {
				for i := 1; i < len(coordinates); i++ {
					var ringArea = area(coordinates[i]...)
					if ringArea < 0 {
						partCoordinates := make([][][]float64, 1)
						partCoordinates[0] = coordinates[i]
						lastAllCoordinates = append(lastAllCoordinates, partCoordinates)
					} else {
						lastAllCoordinates[len(lastAllCoordinates)-1] = append(lastAllCoordinates[len(lastAllCoordinates)-1], coordinates[i])
					}
				}
			}

			//feature = geojson.NewPolygonFeature(coordinates)
			feature = geojson.NewMultiPolygonFeature(lastAllCoordinates...)
		case *shp.MultiPoint:
			mpoint, _ := s.(*shp.MultiPoint)
			coordinates := make([][]float64, mpoint.NumPoints)
			var i int32
			for i = 0; i < mpoint.NumPoints; i = i + 1 {

				coordinates[i] = make([]float64, mpoint.NumPoints)
				cnt := 0
				for j := 0; j < int(mpoint.NumPoints); j = j + 1 {
					coordinates[j] = make([]float64, 2)
					coordinates[j][0] = mpoint.Points[j].X
					coordinates[j][1] = mpoint.Points[j].Y
					cnt = cnt + 1
				}
			}
			feature = geojson.NewMultiPointFeature(coordinates...)

		default:
			fmt.Println("Not support geometry type", reflect.TypeOf(s).Elem())
			continue
		}

		for k, f := range fields {
			val := shape.ReadAttribute(n, k)
			feature.Properties[f.String()] = val
		}
		fc.AddFeature(feature)
	}

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rawJSON, nil
}

func area(points ...[]float64) float64 {
	if len(points) < 3 {
		return 0.0
	}
	var total float64 = 0.0

	for i := 0; i < len(points); i++ {
		xi := points[i][0]
		yi := points[i][1]
		var xi1, yi1 float64
		if i == len(points)-1 {
			xi1 = points[0][0]
			yi1 = points[0][1]

		} else {
			xi1 = points[i+1][0]
			yi1 = points[i+1][1]

		}

		total = total + (xi+xi1)*(yi1-yi)
	}
	return total
}
