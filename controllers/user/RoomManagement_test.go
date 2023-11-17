package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSearching(t *testing.T) {
	test := []struct {
		name         string
		body         models.SearchRequest
		route        string
		errorResult  map[string]string
		expectResult string
	}{{
		name: "from date error",
		body: models.SearchRequest{
			LocOrPlace:       "kannur",
			FromDate:         "02/10/2001",
			ToDate:           "05/10/2001",
			NumberOfChildren: 3,
			NumberOfAdults:   4,
		},
		route:       "/user/home/search",
		errorResult: map[string]string{"error": "{\"error\":\"Invalid from_date format\"}"},
	}, {
		name: "searching success",
		body: models.SearchRequest{
			LocOrPlace:       "Kannur",
			FromDate:         "2010-11-10",
			ToDate:           "2010-11-20",
			NumberOfChildren: 2,
			NumberOfAdults:   3,
		},
		route:        "/user/home/search",
		errorResult:  nil,
		expectResult: "{\"available rooms and counts\":\"hotelDetails\",\"hotels\":\"hotels\"}",
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchAvailableHotels = func(loc string, db *gorm.DB) ([]models.Hotels, error) {
				var hotels []models.Hotels
				return hotels, nil
			}
			fetchAvailableRooms = func(hotelId uint, fromDate, toDate time.Time, NumberOfAdults, NumberOfChildren int, db *gorm.DB) ([]models.Rooms, error) {
				var rooms []models.Rooms
				return rooms, nil
			}

			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}
			r := strings.NewReader(string(body))
			req, err := http.NewRequest(http.MethodPost, tc.route, r)
			if err != nil {
				t.Fatal(err)
			}
			engine.ServeHTTP(w, req)
			if tc.errorResult != nil {
				errResult := tc.errorResult["error"]
				require.Equal(t, w.Body.String(), errResult)
			} else {
				require.Equal(t, w.Code,200)
			}

		})
	}
}

func TestRoomsView(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult string
	}{{
		name:        "RoomView Success",
		route:       "/user/home/rooms?page=1",
		errorResult: "",
	}, {
		name:        "Page number Error",
		route:       "/user/home/rooms?page=0",
		errorResult: "{\"error\":\"Invalid page value\"}",
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchAllRooms = func(skip, limit int, db *gorm.DB) ([]models.Rooms, error) {
				var rooms []models.Rooms
				return rooms, nil
			}

			fetchRoomCategory = func(db *gorm.DB) ([]models.RoomCategory, error) {
				var roomcatagory []models.RoomCategory
				return roomcatagory, nil
			}

			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.route, nil)
			if err != nil {
				require.NoError(t, err)
			}
			engine.ServeHTTP(w, req)
			if tc.errorResult != "" {
				require.Equal(t, w.Body.String(), tc.errorResult)
			} else {
				require.Equal(t, w.Code, 200)
			}
		})
	}
}

func TestRoomDetails(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult string
	}{{
		name:        "View One Room Success",
		route:       "/user/home/rooms/room?id=1",
		errorResult: "",
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchRoomById = func(roomID uint, db *gorm.DB) (*models.Rooms, error) {
				return &rooms, nil
			}
			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.route, nil)
			if err != nil {
				require.NoError(t, err)
			}
			engine.ServeHTTP(w, req)
			if tc.errorResult != "" {
				require.Equal(t, w.Body.String(), tc.errorResult)
			} else {
				require.Equal(t, w.Code, 200)
			}
		})
	}
}
