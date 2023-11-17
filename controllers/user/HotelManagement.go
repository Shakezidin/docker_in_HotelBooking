package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var (
	hotel          = models.Hotels{}
	fetchHotelById = hotel.FetchHotelById
)

// ViewSpecificHotel retrieves a specific hotel by its ID.
func ViewSpecificHotel(c *gin.Context) {
	hotelIDStr := c.Query("id")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotel ID query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "conversion error"})
		return
	}

	_, errr := fetchHotelById(uint(hotelID), Init.DB)
	if errr != nil {
		c.JSON(400, gin.H{"error": "Hotel fetching errror"})
		return
	}
	c.JSON(200, gin.H{"hotel": "hotel"})
}
