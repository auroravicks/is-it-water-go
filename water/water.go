package water

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go-local/is-it-water/utils"
	"io"
	"math"
	"os"

	"github.com/google/tiff"
)

const dataset_path = "./dataset"

type TIFFMetadata struct {
	Width                     int
	Length                    int
	BitsPerSample             int
	Compression               int
	PhotometricInterpretation int
	FillOrder                 int
	StripOffsets              []uint32
	Orientation               int
	SamplesPerPixel           int
	RowsPerStrip              int
	StripByteCounts           []uint32
	XResolution               float64
	YResolution               float64
	PlanarConfiguration       int
	PixelScale                []float64
	ModelTiePoint             []float64
}

func Is_it_water(lat float64, lon float64) string {

	zip_path, tile_id := utils.Zip_path(lat, lon, dataset_path)

	_, err := os.Open(zip_path)
	if err != nil {
		fmt.Println("Could not find file: ", zip_path)
		return "UNKNOWN"
	}

	tif_name := tile_id + "_att.tif"

	r, err := zip.OpenReader(zip_path)
	if err != nil {
		fmt.Println("Could not open file:", zip_path)
		return "UNKNOWN"
	}
	defer r.Close()
	for _, f := range r.File {
		if (f.Name) == tif_name {
			rc, error := f.Open()
			if error != nil {
				fmt.Println("Could not open file: ", tif_name)
				fmt.Println("Error:", error)
				return "UNKNOWN"
			}

			defer rc.Close()
			content, err := io.ReadAll(rc)
			if err != nil {
				fmt.Println("Failed to open file, ", err)
				return "UNKNOWN"
			}

			readerAt := bytes.NewReader(content)

			img, err := tiff.Parse(
				bytes.NewReader(content),
				tiff.DefaultTagSpace,
				tiff.DefaultFieldTypeSpace,
			)
			if err != nil {
				fmt.Println("Dataset object could not be opened: ", err)
				return "UNKNOWN"
			}

			ifd := img.IFDs()[0]

			metadata := extract_metadata(ifd)

			tie_points := metadata.ModelTiePoint
			scale := metadata.PixelScale
			width, height := metadata.Width, metadata.Length

			scaleX, scaleY := scale[0], scale[1]
			originX, originY := tie_points[3], tie_points[4]

			left := originX
			top := originY
			right := left + float64(width-1)*scaleX
			bottom := top - float64(height-1)*scaleY

			bboxWidth := right - left
			bboxHeight := top - bottom

			widthPct := (lon - left) / bboxWidth
			heightPct := (lat - bottom) / bboxHeight

			x := int(math.Floor(float64(width) * widthPct))
			y := int(math.Floor(float64(height) * (1 - heightPct)))

			if x < 0 || x >= width || y < 0 || y >= height {
				fmt.Print("\n Exceptional X and Y values")
				return "UNKNOWN"
			}

			stripIndex := int(math.Floor(float64(y) / float64(metadata.RowsPerStrip)))
			offset := metadata.StripOffsets[stripIndex]
			byteCount := metadata.StripByteCounts[stripIndex]

			sectionReader := io.NewSectionReader(readerAt, int64(offset), int64(byteCount))
			stripData := make([]byte, byteCount)
			_, err = sectionReader.Read(stripData)
			if err != nil {
				fmt.Println("Error while opening up the file")
				return "UNKNOWN"
			}

			rowStart := (y % metadata.RowsPerStrip) * width
			pixelIndex := rowStart + x
			pixelValue := stripData[pixelIndex]

			switch pixelValue {
			case 0:
				return "LAND"
			case 1:
				return "OCEAN"
			case 2:
				return "RIVER"
			case 3:
				return "LAKE"
			default:
				return "UNKNOWN"
			}
		}
	}
	return "UNKNOWN"
}

// func main() {
// 	fmt.Println("\n", (is_it_water(0, 0)))
// 	fmt.Println("\n", (is_it_water(9999, 9999)))
// 	fmt.Println("\n", (is_it_water(46.67, 103.3)))
// 	fmt.Println("\n", (is_it_water(-10.47, 105.57)))
// 	fmt.Println("\n", (is_it_water(30.03, 31.22)))
// 	fmt.Println("\n", (is_it_water(45.295504, 12.61337)))
// 	fmt.Println("\n", (is_it_water(44.35953, -2.765543)))
// 	fmt.Println("\n", (is_it_water(6.56, 0)))
// 	fmt.Println("\n", (is_it_water(0, -69)))
// }

func extract_metadata(ifd tiff.IFD) TIFFMetadata {
	metadata := TIFFMetadata{}

	metadata.Width = utils.SingleInt(ifd.GetField(256))
	metadata.Length = utils.SingleInt(ifd.GetField(257))
	metadata.BitsPerSample = utils.SingleInt(ifd.GetField(258))
	metadata.Compression = utils.SingleInt(ifd.GetField(259))
	metadata.PhotometricInterpretation = utils.SingleInt(ifd.GetField(262))
	metadata.FillOrder = utils.SingleInt(ifd.GetField(266))
	metadata.StripOffsets = utils.LongArray(ifd.GetField(273))
	metadata.Orientation = utils.SingleInt(ifd.GetField(274))
	metadata.SamplesPerPixel = utils.SingleInt(ifd.GetField(277))
	metadata.RowsPerStrip = utils.SingleInt(ifd.GetField(278))
	metadata.StripByteCounts = utils.LongArray(ifd.GetField(279))
	metadata.XResolution = utils.Rational(ifd.GetField(282))
	metadata.YResolution = utils.Rational(ifd.GetField(283))
	metadata.PlanarConfiguration = utils.SingleInt(ifd.GetField(284))

	metadata.PixelScale = utils.FloatArray(ifd.GetField(33550))
	metadata.ModelTiePoint = utils.FloatArray(ifd.GetField(33922))

	return metadata
}
