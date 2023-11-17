package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var createcontact = models.Createcontact

// SubmitContact handles submitting a contact message to the admin.
func SubmitContact(c *gin.Context) {
	var message = models.Message{}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding error"})
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	user, err := fetchUser(username, Init.DB)
	if err != nil {
		c.JSON(400, gin.H{"Error": "User Fetching error"})
		return
	}

	contact := &models.Contact{
		Message: message.Message,
		UserID:  user.ID,
	}

	errr := createcontact(contact, Init.DB)
	if errr != nil {
		c.JSON(400, gin.H{"Error": "Message creation error"})
		return
	}

	c.Status(http.StatusOK)
}
