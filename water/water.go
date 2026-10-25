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

const defaultDatasetPath = "./dataset"

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
	datasetPath := os.Getenv("DATASET_PATH")
	if datasetPath == "" {
		datasetPath = defaultDatasetPath
	}
	zip_path, tile_id := utils.Zip_path(lat, lon, datasetPath)

	zipFile, err := os.Open(zip_path)
	if err != nil {
		fmt.Println("Could not find file: ", zip_path)
		return "UNKNOWN"
	}
	defer zipFile.Close()

	fileInfo, err := zipFile.Stat()
	if err != nil {
		fmt.Println("Could not get file info:", err)
		return "UNKNOWN"
	}

	r, err := zip.NewReader(zipFile, fileInfo.Size())
	if err != nil {
		fmt.Println("Could not open zip archive:", zip_path)
		return "UNKNOWN"
	}

	tif_name := tile_id + "_att.tif"
	var tifFile *zip.File
	for _, f := range r.File {
		if f.Name == tif_name {
			tifFile = f
			break
		}
	}

	if tifFile == nil {
		fmt.Println("Could not find TIFF file: ", tif_name)
		return "UNKNOWN"
	}

	rc, err := tifFile.Open()
	if err != nil {
		fmt.Println("Could not open TIFF file:", tif_name)
		return "UNKNOWN"
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		fmt.Println("Failed to read TIFF file:", err)
		return "UNKNOWN"
	}

	readerAt := bytes.NewReader(content)
	img, err := tiff.Parse(readerAt, tiff.DefaultTagSpace, tiff.DefaultFieldTypeSpace)
	if err != nil {
		fmt.Println("Error parsing TIFF file:", err)
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
	right := left + float64(width)*scaleX
	bottom := top - float64(height)*scaleY

	bboxWidth := right - left
	bboxHeight := top - bottom

	widthPct := (lon - left) / bboxWidth
	heightPct := (lat - bottom) / bboxHeight

	x := int(math.Floor(float64(width) * widthPct))
	y := int(math.Floor(float64(height) * (1 - heightPct)))

	fmt.Println(x, y)
	if x < 0 || x >= width || y < 0 || y >= height {
		fmt.Println("Exceptional X and Y values")
		return "UNKNOWN"
	} else if x == width {
		x -= 1
	} else if y == height {
		y -= 1
	}

	stripIndex := int(math.Floor(float64(y) / float64(metadata.RowsPerStrip)))
	offset := metadata.StripOffsets[stripIndex]
	byteCount := metadata.StripByteCounts[stripIndex]

	sectionReader := io.NewSectionReader(readerAt, int64(offset), int64(byteCount))
	stripData := make([]byte, byteCount)
	_, err = sectionReader.Read(stripData)
	if err != nil {
		fmt.Println("Error parsing TIFF file: ", err)
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
