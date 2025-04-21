package main

import (
	"go-local/is-it-water/water"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/is-it-water/:lat/:lon", func(c *gin.Context) {
		start := time.Now()

		lat, _ := strconv.ParseFloat(c.Param("lat"), 64)
		lon, _ := strconv.ParseFloat(c.Param("lon"), 64)

		feature := water.Is_it_water(lat, lon)
		c.JSON(200, gin.H{
			"isWater": feature != "LAND",
			"feature": feature,
			"lat":     lat,
			"lon":     lon,
			"reqMs":   time.Since(start).Milliseconds(),
		})
	})
	router.Run("127.0.0.1:8000")
}
