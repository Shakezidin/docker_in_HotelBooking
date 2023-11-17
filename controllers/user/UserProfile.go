package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var (
	updateUser     = models.UpdateUSer
	userr          = &models.User{}
	checkPasswordd = userr.CheckPassword
	booking        = models.Booking{}
	history        = booking.FetchHistory
)

// Profile handles user profile retrieval.
func Profile(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")

	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}

	user, errr := fetchUser(username, Init.DB)
	if errr != nil {
		c.JSON(400, gin.H{"Error": "Fetching user error"})
		return
	}

	c.JSON(200, gin.H{"Status": "Success","User":user})
}

// ProfileEdit handles editing user profile.
func ProfileEdit(c *gin.Context) {

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	user, err := fetchUser(username, Init.DB)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Fetching user error"})
		return
	}

	var updateuser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := c.BindJSON(&updateuser); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	if updateuser.Email == "" {
		updateuser.Email = user.Email
	}

	if updateuser.Phone == "" {
		updateuser.Phone = user.Phone
	}

	if updateuser.Name == "" {
		updateuser.Name = user.Name
	}

	user.Name = updateuser.Name
	user.Email = updateuser.Email
	user.Phone = updateuser.Phone

	if err := updateUser(user, Init.DB); err != nil {
		c.JSON(400, gin.H{"Error": "Error while updating user"})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

// PasswordChange handles changing user password.
func PasswordChange(c *gin.Context) {
	var pswrd struct {
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}

	if err := c.BindJSON(&pswrd); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}

	user, err := fetchUser(username, Init.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found"})
		return
	}
	if err := checkPasswordd(pswrd.OldPassword); err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}

	if err := user.HashPassword(pswrd.NewPassword); err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing error",
		})
		c.Abort()
		return
	}

	err = updateUser(user, Init.DB)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Updating user error"})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

// History handles user booking history retrieval.
func History(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}

	user, err := fetchUser(username, Init.DB)
	if err != nil {
		c.JSON(400, gin.H{"Error": "User fetching error"})
		return
	}

	_, errr := history(user.ID, Init.DB)
	if errr != nil {
		c.JSON(400, gin.H{"Error": "error while fetching booking"})
		return
	}

	c.JSON(200, gin.H{"history": "booking"})
}
