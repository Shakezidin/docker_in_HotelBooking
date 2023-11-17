package routes

import (
	"github.com/gin-gonic/gin"
	UserCtrl "github.com/shaikhzidhin/controllers/user"
)

// UserRoutes Set up the routes for the user section of the application.
func UserRoutes(c *gin.Engine) {
	UserCtrl.RegisterUserRoutes(c)

}
