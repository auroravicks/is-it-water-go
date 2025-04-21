package utils

import (
	"fmt"
	"math"

	"github.com/google/tiff"
)

func Zip_path(lat, lon float64, dataset_path string) (string, string) {
	lat_prefix, lon_prefix := "N", "E"
	if lat < 0 {
		lat_prefix = "S"
	}
	if lon < 0 {
		lon_prefix = "W"
	}

	lat_int := int(math.Abs(math.Floor(lat))) % 100
	lon_int := int(math.Abs(math.Floor(lon))) % 1000
	lat_str := fmt.Sprintf("%s%02d", lat_prefix, lat_int)
	lon_str := fmt.Sprintf("%s%03d", lon_prefix, lon_int)
	tile_id := "ASTWBDV001_" + lat_str + lon_str
	zip_path := dataset_path + "/" + tile_id + ".zip"

	return zip_path, tile_id
}

func SingleInt(field tiff.Field) int {
	raw := field.Value().Bytes()
	bo := field.Value().Order()

	return int(bo.Uint16(raw[0*4 : (1)*4]))
}

func FloatArray(field tiff.Field) []float64 {
	raw := field.Value().Bytes()
	order := field.Value().Order()
	floats := make([]float64, len(raw)/8)

	for i := range floats {
		bits := order.Uint64(raw[i*8 : (i+1)*8])
		floats[i] = (math.Round((math.Float64frombits(bits)) * 1e6)) / (1e6)
	}
	return floats
}

func LongArray(field tiff.Field) []uint32 {
	raw := field.Value().Bytes()
	order := field.Value().Order()
	count := len(raw) / 4
	longs := make([]uint32, count)

	for i := 0; i < count; i++ {
		longs[i] = order.Uint32(raw[i*4 : (i+1)*4])
	}
	return longs
}

func Rational(field tiff.Field) float64 {
	raw := field.Value().Bytes()
	order := field.Value().Order()

	numerator := order.Uint32(raw[0:4])
	denonimator := order.Uint32((raw[4:8]))

	return float64(numerator) / float64(denonimator)
}
