package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestViewSpecificHotel(t *testing.T) {
	test := []struct {
		name         string
		route        string
		errorResult  map[string]string
		expectResult string
	}{{
		name:         "Success",
		route:        "/user/home/banner/hotel?id=1",
		errorResult:  nil,
		expectResult: "{\"hotel\":\"hotel\"}",
	}}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchHotelById = func(hotelId uint, db *gorm.DB) (*models.Hotels, error) {
				return &hotel, nil
			}
			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			req,err:=http.NewRequest(http.MethodGet,tc.route,nil)
			if err!=nil{
				t.Fatal(err)
			}
			w:=httptest.NewRecorder()
			engine.ServeHTTP(w,req)
			if tc.errorResult!=nil{
				errValue := tc.errorResult["Error"]
				require.Equal(t,w.Body.String(),errValue)
			}else{
				require.Equal(t,w.Body.String(),tc.expectResult)
			}

		})
	}
}
