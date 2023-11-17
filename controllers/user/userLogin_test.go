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

func TestLogin(t *testing.T) {
	// Create a new instance of the MockDB
	test := []struct {
		name        string
		body        models.Login
		route       string
		errorResult map[string]string
	}{
		{
			name: "Success",
			body: models.Login{
				Username: "Shaikh_Zidhin",
				Password: "Sinu1090.",
			},
			route:       "/user/login",
			errorResult: nil,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchUser=func(val string, db *gorm.DB) (*models.User, error) {
				return &user,nil
			}
			checkPassword = func(providedPassword string) error {
				return nil
			}

			body, err := json.Marshal(tc.body)
			if err != nil {
				require.NoError(t, err)
			}

			r := strings.NewReader(string(body))

			w, err := Setup(http.MethodPost, tc.route, r,"")
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
