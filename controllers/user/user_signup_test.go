package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		name        string
		body        models.User
		route       string
		errorResult map[string]string
	}{

		{
			name: "error- binding_error",
			body: models.User{
				UserName: "",
				Name:     "test_name",
				Email:    "test_email",
				Phone:    "75217504332",
				Password: "Sinu1090.",
			},
			route:       "/user/signup",
			errorResult: map[string]string{"error": "validation error1"},
		},
		{
			name: "success",
			body: models.User{
				UserName: "Sinu_zidin",
				Name:     "test_name",
				Email:    "test_email",
				Phone:    "75217504332",
				Password: "Sinu1090.",
			},
			route:       "/user/signup",
			errorResult: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getOtp = func(name, email string) string {
				return "1234"
			}
			setRedis = func(key string, value any, expirationTime time.Duration) error {
				return nil
			}

			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}
			r := strings.NewReader(string(body))

			w, err := Setup(http.MethodPost, tc.route, r, "")
			if err != nil {
				require.NoError(t, err)

			}
			if tc.errorResult != nil {
				errValue, _ := json.Marshal(tc.errorResult)
				// require.NoError(t, err)
				require.JSONEq(t, w.Body.String(), string(errValue))
			} else {
				data, err := readJson("testdata/user_signup.json")
				if err != nil {
					require.NoError(t, err)
				}

				require.JSONEq(t, w.Body.String(), data)
			}

		})

	}

}

func Setup(method, url string, body io.Reader, token string) (*httptest.ResponseRecorder, error) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	engine := gin.Default()

	RegisterUserRoutes(engine)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	engine.ServeHTTP(w, req)
	return w, nil
}

func readJson(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("error opening file")
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestSignupVerification(t *testing.T) {
	test := []struct {
		name        string
		body        models.OtpCredentials
		route       string
		errorResult map[string]string
	}{
		{
			name: "Success",
			body: models.OtpCredentials{
				Email: "sinuzidin@gmail.com",
				Otp:   "1234",
			},
			route:       "/user/signup/verification",
			errorResult: nil,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the necessary functions
			verifyOtp = func(superkey, otpInput string, c *gin.Context) bool {
				return true
			}
			getRedis = func(key string) (string, error) {
				jsonData, err := json.Marshal(tc.body)
				if err != nil {
					return "", err
				}
				return string(jsonData), nil
			}
			create = func(userr *models.User, db *gorm.DB) error {
				return nil
			}
			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}
			r := strings.NewReader(string(body))

			// Simulate an HTTP request to your function
			w, err := Setup(http.MethodPost, tc.route, r, "")
			if err != nil {
				require.NoError(t, err)
			}

			if tc.errorResult != nil {
				// Check for an error response
				errValue, _ := json.Marshal(tc.errorResult)
				require.JSONEq(t, w.Body.String(), string(errValue))
			} else {
				// Check for a success response
				data, err := readJson("testdata/user_singnup_success.json")
				if err != nil {
					require.NoError(t, err)
				}
				require.JSONEq(t, w.Body.String(), data)
			}
		})
	}
}
