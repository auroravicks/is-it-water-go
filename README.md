# Is It Water? (Go Implementation)

A high-performance Go implementation of the [Is It Water?](https://is-on-water.balbona.me/) service. This service determines whether a given geographic coordinate (latitude/longitude) is located on water or land, using NASA's ASTWBD dataset.

## Features

- Fast water/land detection using pre-processed geospatial data
- Simple HTTP API endpoint
- Low memory footprint
- Efficient tile-based data loading
- Built with Go for high performance

## Prerequisites

- Go 1.19 or later
- Sufficient disk space for the dataset (approximately 100MB)

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/is-it-water-go.git
   cd is-it-water-go
   ```

2. Download the required dataset:
   ```bash
   # Create dataset directory
   mkdir -p dataset
   
   # Download the dataset (example for a single tile, you'll need all required tiles)
   # Note: You'll need to obtain the ASTWBD dataset tiles and place them in the dataset/ directory
   # The expected format is ASTWBDV001_N00E000.zip for each tile
   ```

### Configuration

The service is pre-configured to look for dataset files in the `./dataset` directory. You can modify this by changing the `dataset_path` constant in `water/water.go`.

### Running the Service

1. Start the server:
   ```bash
   go run main.go
   ```

2. The service will be available at `http://localhost:8000`

## API Usage

### Check if a point is on water

```
GET /is-it-water/{lat}/{lon}
```

#### Parameters

- `lat`: Latitude (decimal degrees, -90 to 90)
- `lon`: Longitude (decimal degrees, -180 to 180)

#### Example Request

```
GET http://localhost:8000/is-it-water/40.7128/-74.0060
```

#### Example Response

```json
{
  "isWater": false,
  "feature": "LAND",
  "lat": 40.7128,
  "lon": -74.0060,
  "reqMs": 12
}
```

#### Response Fields

- `isWater`: Boolean indicating if the point is water
- `feature`: String describing the feature type (e.g., "LAND", "WATER")
- `lat`: The latitude provided in the request
- `lon`: The longitude provided in the request
- `reqMs`: Request processing time in milliseconds

## Dataset

This project uses the ASTWBD (ASTER Water Body Dataset) from NASA, which provides global water body data at 1 arc-second (approximately 30m) resolution. The dataset is already included in the repository under the `dataset/` directory, containing multiple ZIP files where each file covers a 1°×1° tile.
