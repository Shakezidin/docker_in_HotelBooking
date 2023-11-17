package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(c *gin.Engine) {
	user := c.Group("/user")

	// User Authentication routes
	user.POST("/login", Login)
	user.POST("/signup", Signup)
	user.POST("/signup/verification", SignupVerification)

	// Password Recovery routes
	user.POST("/password/forget", ForgetPassword)
	user.POST("/password/forget/verifyotp", VerifyOTP)
	user.POST("/password/set/new", NewPassword)

	// User Home & Profile routes
	// user.Use(middleware.UserAuthMiddleware)
	user.GET("/", Home)
	user.GET("/profile", Profile)
	user.PATCH("/profile/edit", ProfileEdit)
	user.PUT("/profile/password/change", PasswordChange)
	user.GET("/booking/history", History)

	// Hotels routes
	user.GET("/home", Home)
	user.GET("/home/banner", BannerShowing)
	user.GET("/home/banner/hotel", ViewSpecificHotel)
	user.POST("/home/search", Searching)

	// Room routes
	user.GET("/home/rooms", RoomsView)
	user.GET("/home/rooms/room", RoomDetails)
	user.POST("/home/rooms/filter", RoomFilter)

	// Contact routes
	user.POST("/home/contact", SubmitContact)

	// Booking Management routes
	user.GET("/home/room/book", CalculateAmountForDays)
	user.GET("/coupons/view", ViewNonBlockedCoupons)
	user.GET("/coupon/apply", ApplyCoupon)
	user.GET("/wallet", ViewWallet)
	user.GET("/wallet/apply", ApplyWallet)
	user.GET("/payat/hotel", OfflinePayment)

	// Razorpay routes
	user.GET("/online/payment", RazorpayPaymentGateway)
	user.GET("/payment/success", RazorpaySuccess)
	user.GET("/success", SuccessPage)
	user.GET("/cancel/booking", CancelBooking)
}
