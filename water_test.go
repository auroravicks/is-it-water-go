package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const BASE_URL = "http://localhost:8000/is-it-water"

type ApiResponse struct {
	IsWater bool    `json:"isWater"`
	Feature string  `json:"feature"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	ReqMs   int64   `json:"reqMs"`
}

func TestIsInWater(t *testing.T) {
	tests := []struct {
		lat             float64
		lon             float64
		expectedFeature string
		expectedIsWater bool
	}{
		{0, 0, "UNKNOWN", true},
		{9999, 9999, "UNKNOWN", true},
		{46.67, 103.3, "LAKE", true},
		{-10.47, 105.57, "LAND", false},
		{30.03, 31.22, "RIVER", true},
		{45.295504, 12.61337, "OCEAN", true},
		{44.35953, -2.765543, "UNKNOWN", true},
		{6.56, 0, "LAKE", true},
		{0, -69, "LAND", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing lat: %f, lon: %f", tt.lat, tt.lon), func(t *testing.T) {
			// Make the HTTP GET request
			resp, err := http.Get(fmt.Sprintf("%s/%f/%f", BASE_URL, tt.lat, tt.lon))
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check the response status code
			assert.Equal(t, 200, resp.StatusCode)

			var data ApiResponse
			// Decode the JSON response
			err = json.NewDecoder(resp.Body).Decode(&data)
			assert.NoError(t, err)

			// Print the values for debugging
			fmt.Printf("\nTesting (%f, %f)\n", tt.lat, tt.lon)
			fmt.Printf("→ Expected Feature: %s, Actual Feature: %s\n", tt.expectedFeature, data.Feature)
			fmt.Printf("→ Expected isWater: %t, Actual isWater: %t\n", tt.expectedIsWater, data.IsWater)

			// Assertions
			assert.Equal(t, tt.expectedFeature, data.Feature)
			assert.Equal(t, tt.expectedIsWater, data.IsWater)
		})
	}
}
