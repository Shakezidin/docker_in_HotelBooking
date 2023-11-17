package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var authToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFhcm9vIiwiZXhwIjoxNzAwNDkxMzY2fQ.sKpi_mbDJ3QvfYl8j7n7qzPRooE2pkaPEK5aFWowaww"

func TestProfile(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
	}{{
		name:        "Successs",
		route:       "/user/profile",
		errorResult: nil,
	}}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchUser = func(val string, db *gorm.DB) (*models.User, error) {
				return &user, nil
			}

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			engine := gin.Default()

			RegisterUserRoutes(engine)
			req, err := http.NewRequest(http.MethodGet,tc.route, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", authToken)
			engine.ServeHTTP(w, req)
			if tc.errorResult != nil {
				errValue, _ := json.Marshal(tc.errorResult)
				require.JSONEq(t, w.Body.String(), string(errValue))
			} else {
				require.Equal(t, w.Code,200)
			}

		})
	}
}

func TestProfileEdit(t *testing.T) {

	type updateuser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	test := []struct {
		name        string
		body        updateuser
		route       string
		errorResult map[string]string
	}{{
		name: "Successs",
		body: updateuser{
			Name:  "Shaikh",
			Email: "Sinuzidin@gmail.com",
			Phone: "9061978992",
		},
		route:       "/user/profile/edit",
		errorResult: nil,
	}}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchUser = func(val string, db *gorm.DB) (*models.User, error) {
				return &user, nil
			}
			updateUser = func(usermodel *models.User, db *gorm.DB) error {
				return nil
			}
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			engine := gin.Default()
			RegisterUserRoutes(engine)
			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}
			r := strings.NewReader(string(body))
			req, err := http.NewRequest(http.MethodPatch, tc.route, r)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", authToken)
			engine.ServeHTTP(w, req)
			if tc.errorResult == nil {
				data, err := readJson("testdata/user_profEdit.json")
				if err != nil {
					require.NoError(t, err)
				}
				require.JSONEq(t, w.Body.String(), data)
			}

		})
	}
}

func TestPasswordChange(t *testing.T) {
	type pswrd struct {
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}
	test := []struct {
		name         string
		body         pswrd
		route        string
		errorResult  map[string]string
		expectResult string
	}{{
		name: "Success",
		body: pswrd{
			OldPassword: "Sinu1090.",
			NewPassword: "Sinu1090#",
		},
		route:        "/user/profile/password/change",
		errorResult:  nil,
		expectResult: "{\"status\":\"success\"}",
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchUser = func(val string, db *gorm.DB) (*models.User, error) {
				return &user, nil
			}
			updateUser = func(user *models.User, db *gorm.DB) error {
				return nil
			}
			checkPasswordd = func(providedPassword string) error {
				return nil
			}
			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			w := httptest.NewRecorder()
			RegisterUserRoutes(engine)
			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}
			bod := strings.NewReader(string(body))
			req, err := http.NewRequest(http.MethodPut, tc.route, bod)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", authToken)
			engine.ServeHTTP(w, req)
			if tc.errorResult != nil {
				errValue := tc.errorResult["Error"]
				require.Equal(t, w.Body.String(), errValue)
			} else {
				require.Equal(t, w.Body.String(), tc.expectResult)
			}
		})
	}
}

func TestHistory(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
		expectResult string
	}{{
		name: "Success",
		route: "/user/booking/history",
		errorResult: nil,
		expectResult:"{\"history\":\"booking\"}",
	}}

	for _,tc:=range test{
		t.Run(tc.name,func(t *testing.T) {
			fetchUser=func(val string, db *gorm.DB) (*models.User, error) {
				return &user,nil
			}
			history=func(userId uint, db *gorm.DB) (*models.Booking, error) {
				return &booking,nil
			}
			gin.SetMode(gin.TestMode)
			engine:=gin.Default()
			RegisterUserRoutes(engine)
			w:=httptest.NewRecorder()
			req,err:=http.NewRequest(http.MethodGet,tc.route,nil)
			if err!=nil{
				require.NoError(t,err)
			}
			req.Header.Set("Authorization",authToken)
			engine.ServeHTTP(w,req)
			if tc.errorResult!=nil{
				errResult:=tc.errorResult["Error"]
				require.Equal(t,w.Body.String(),errResult)
			}else{
				require.Equal(t,w.Body.String(),tc.expectResult)
			}
		})
	}
}
