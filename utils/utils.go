package utils

import (
	"fmt"
	"math"
	"math/big"

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
		floats[i] = math.Float64frombits(bits)
	}

	return floats
}

func Rational(field tiff.Field) float64 {
	raw := field.Value().Bytes()
	order := field.Value().Order()

	numerator := order.Uint32(raw[0:4])
	denominator := order.Uint32((raw[4:8]))

	bfNum := new(big.Float).SetInt64(int64(numerator))
	bfDen := new(big.Float).SetInt64(int64(denominator))

	result := new(big.Float).Quo(bfNum, bfDen)

	result = result.SetPrec(256)
	result = result.SetMode(big.ToNearestEven)

	value, _ := result.Float64()
	return value
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
