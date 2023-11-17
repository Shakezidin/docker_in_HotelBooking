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

func TestSubmitContact(t *testing.T) {
	test := []struct {
		name        string
		route       string
		body        models.Message
		errorResult string
	}{{
		name:  "Contact Submit Success",
		route: "/user/home/contact",
		body: models.Message{
			Message: "helloooo",
		},
		errorResult: "",
	}}
	for _, tc := range test {
		fetchUser = func(val string, db *gorm.DB) (*models.User, error) {
			return &user, nil
		}
		createcontact = func(contact *models.Contact, db *gorm.DB) error {
			return nil
		}

		gin.SetMode(gin.TestMode)
		engine := gin.Default()
		RegisterUserRoutes(engine)
		w := httptest.NewRecorder()
		body,err:=json.Marshal(tc.body)
		if err!=nil{
			require.NoError(t,err)
		}
		r:=strings.NewReader(string(body))
		req, err := http.NewRequest(http.MethodPost, tc.route, r)
		if err != nil {
			require.NoError(t, err)
		}
		req.Header.Set("Authorization", authToken)
		engine.ServeHTTP(w, req)
		if tc.errorResult != "" {
			errValue, _ := json.Marshal(tc.errorResult)
			require.JSONEq(t, w.Body.String(), string(errValue))
		} else {
			require.Equal(t, w.Code, 200)
		}
	}
}
