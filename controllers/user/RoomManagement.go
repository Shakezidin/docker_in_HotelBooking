package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var (
	hotelref             = models.Hotels{}
	roomref              = models.Rooms{}
	categoryref          = models.RoomCategory{}
	fetchAvailableHotels = hotelref.FetchAvailableHotels
	fetchAvailableRooms  = roomref.FetchingAvailableRooms
	fetchAllRooms        = roomref.FetchAllRooms
	fetchRoomCategory    = categoryref.FetchRoomCategory
	fetchRoomById        = roomref.FetchRoomById
)

// Searching helps to find available rooms and hotels in a loacation
func Searching(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	fromDate, err := time.Parse(layout, req.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_date format"})
		return
	}
	toDate, err := time.Parse(layout, req.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_date format"})
		return
	}

	err = setRedis("fromdate", req.FromDate, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting 'fromdate' in Redis client"})
		return
	}

	err = setRedis("todate", req.ToDate, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting 'todate' in Redis client"})
		return
	}

	// Fetch hotels that match the location or place
	hotels, err := fetchAvailableHotels(req.LocOrPlace, initializer.DB)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Fetching available hotel eror"})
		return
	}

	// Create a list to store hotel details including room category details
	hotelDetails := make([]map[string]interface{}, 0)

	for _, hotel := range hotels {
		// Fetch available rooms for the hotel
		availableRooms, err := fetchAvailableRooms(hotel.ID, fromDate, toDate, req.NumberOfAdults, req.NumberOfChildren, initializer.DB)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching rooms for the hotel"})
			return
		}

		// Create a list to store room details for each category
		categoryDetails := make(map[string][]map[string]interface{})
		addedCategories := make(map[string]bool)

		for _, room := range availableRooms {
			category := room.RoomCategory.Name
			if categoryDetails[category] == nil {
				categoryDetails[category] = make([]map[string]interface{}, 0)
			}

			// Add room details to the category only if it hasn't been added before
			if !addedCategories[category] {
				roomDetails := map[string]interface{}{
					"room_id":     room.ID,
					"description": room.Description,
					"price":       room.Price,
					"adults":      room.Adults,
					"children":    room.Children,
					"bed":         room.Bed,
					"images":      room.Images,
				}
				categoryDetails[category] = append(categoryDetails[category], roomDetails)
				addedCategories[category] = true
			}
		}

		// Calculate the available room count for each category
		categoryCounts := make(map[string]int)
		for _, room := range availableRooms {
			categoryCounts[room.RoomCategory.Name]++
		}

		// Add hotel details including room category details and room counts to the list
		hotelDetails = append(hotelDetails, map[string]interface{}{
			"hotel_name":           hotel.Name,
			"place":                hotel.City,
			"facilities":           hotel.Facility,
			"category_details":     categoryDetails,
			"available_room_count": categoryCounts, // Add available room counts
		})
	}

	c.JSON(http.StatusOK, gin.H{"hotels": "hotels", "available rooms and counts": hotelDetails})
}

// RoomsView returns a list of rooms for viewing.
func RoomsView(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page value"})
		return
	}
	limit := 10
	skip := (pageInt - 1) * limit

	rooms, errr := fetchAllRooms(skip, limit, initializer.DB)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	categories, errrr := fetchRoomCategory(initializer.DB)
	if errrr != nil {
		c.JSON(400, gin.H{"error": "Catagory fetching error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms, "categories": categories})
}

// RoomDetails returns details of a specific room.
func RoomDetails(c *gin.Context) {
	roomIDStr := c.Query("id")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	room, errr := fetchRoomById(uint(roomID), initializer.DB)
	if errr != nil {
		c.JSON(400, gin.H{"Error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}
