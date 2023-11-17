package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var (
	banner      = models.Banner{}
	hotels      = models.Hotels{}
	rooms       = models.Rooms{}
	fetchbanner = banner.FetchBanner
	fetchHotels = hotels.FetchHotels
	fetchRooms  = rooms.FetchinRooms
)

// Home is a handler for the user's homepage.
func Home(c *gin.Context) {
	banners, err := fetchbanner(true, true, Init.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while fetching banners"})
		return
	}
	city := c.Query("loc")
	if city == "" {
		c.JSON(400, gin.H{"error": "location query parameter is missing"})
		return
	}

	page := c.Query("page")
	limit := 10
	pageInt, _ := strconv.Atoi(page)

	skip := (pageInt - 1) * limit

	var rooms []models.Rooms

	// Retrieve hotels based on the location and pagination
	hotels, err := fetchHotels(city, true, false, true, skip, limit, Init.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Fetching hotel error"})
		return
	}

	// Retrieve rooms for each hotel
	for _, hotel := range hotels {

		// Retrieve rooms for the current hotel
		hotelRooms, err := fetchRooms(hotel.ID,Init.DB)
		if err != nil {
			c.JSON(400, gin.H{"Error": "Error fetching rooms"})
			return
		}

		// Append the rooms of the current hotel to the rooms slice
		rooms = append(rooms, hotelRooms...)
	}
	
	c.JSON(200, gin.H{"Hotels": hotels, "Rooms": rooms, "Banners": banners})
}

// BannerShowing is a handler for displaying all banners.
func BannerShowing(c *gin.Context) {
	var banners []models.Banner

	if err := Init.DB.Preload("Hotels").Where("available = ? AND active = ?", true, true).Find(&banners).Error; err != nil {
		c.JSON(400, gin.H{"error": "banner retrieval error"})
		return
	}

	c.JSON(200, gin.H{"status": banners})
}
