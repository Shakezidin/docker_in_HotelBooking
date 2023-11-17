package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/shaikhzidhin/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHome(t *testing.T) {
	test := []struct {
		name        string
		route       string
		body        string
		errorResult map[string]string
	}{{
		name:        "Success",
		route:       "/user/?loc=kannur&page=1",
		body:        "",
		errorResult: nil,
	}}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchbanner = func(available, active bool, db *gorm.DB) ([]*models.Banner, error) {
				var banners []*models.Banner
				return banners, nil
			}
			fetchHotels = func(city string, isAvailab, isBlock, adminApproval bool, skip, limit int, db *gorm.DB) ([]models.Hotels, error) {
				var hotels []models.Hotels
				return hotels, nil
			}

			fetchRooms = func(hotelId uint,db *gorm.DB) ([]models.Rooms, error) {
				var rooms []models.Rooms
				return rooms, nil
			}

			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}

			r := strings.NewReader(string(body))

			w, err := Setup(http.MethodGet, tc.route, r,"")
			if err != nil {
				require.NoError(t, err)
			}
			if tc.errorResult != nil {
				errValue, _ := json.Marshal(tc.errorResult)
				require.JSONEq(t, w.Body.String(), string(errValue))
			} else {
				require.Equal(t, w.Code, 200)
			}
		})
	}
}
