package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var (
	Couponref       = models.Coupon{}
	usedcouponref   = models.UsedCoupon{}
	walletref       = models.Wallet{}
	transaction     = models.Transaction{}
	fetchAllCoupon  = Couponref.FetchAllCoupon
	fetchCouponById = Couponref.FetchCouponById
	fetchUsedCoupon = usedcouponref.FetchUsedCoupon
	fetchWallet     = walletref.FetchWallet
	fetchWalletByID = walletref.FetchWalletById
	saveWallet      = walletref.SaveWallet
	createwallet          = transaction.Create
)

// CalculateAmountForDays calculates the amount for booking based on selected dates and room.
func CalculateAmountForDays(c *gin.Context) {
	roomIDStr := c.Query("id")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "convert error"})
		return
	}

	err = setRedis("roomid", roomIDStr, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	fromdateStr, err := getRedis("fromdate")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'fromdate' from Redis client"})
		return
	}

	todateStr, err := getRedis("todate")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'todate' from Redis client"})
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromdateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromdate format"})
		return
	}

	toDate, err := time.Parse("2006-01-02", todateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toDate format"})
		return
	}

	duration := toDate.Sub(fromDate)
	days := int(duration.Hours() / 24)

	room, err := fetchRoomById(uint(roomID), Init.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetching rooms"})
		return
	}

	totalPrice := days * int(room.DiscountPrice)

	GSTPercentage := 18.0 // Use a floating-point number for the percentage
	GSTAmount := (GSTPercentage / 100.0) * float64(totalPrice)
	payableAmount := totalPrice + int(GSTAmount)

	err = setRedis("Amount", payableAmount, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"roomPrice":     room.Price,
		"discount":      room.Discount,
		"TotalAmount":   totalPrice,
		"GSTAmount":     int(GSTAmount),
		"PayableAmount": payableAmount,
	})
}

// ViewNonBlockedCoupons retrieves non-blocked coupons.
func ViewNonBlockedCoupons(c *gin.Context) {
	var coupons []models.Coupon

	coupons, err := fetchAllCoupon(Init.DB)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Coupon fetching error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

// ApplyCoupon applies a coupon to the booking.
func ApplyCoupon(c *gin.Context) {
	couponID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{"Error": "query id error"})
		return
	}
	amountStr, err := getRedis("Amount")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'Amount' from Redis client: " + err.Error()})
		return
	}

	amount, errrr := strconv.Atoi(amountStr)
	if errrr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
		return
	}

	CouponIDstr, _ := getRedis("couponID")
	if CouponIDstr != "" {
		oldcouponID, err := strconv.Atoi(CouponIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
			return
		}

		if oldcouponID == couponID {
			c.JSON(400, gin.H{"error": "alredy used coupon"})
			return
		}
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}

	user, errr := fetchUser(username, Init.DB)
	if errr != nil {
		c.JSON(400, gin.H{"Error": "User not found"})
		return
	}

	coupon, errrrr := fetchCouponById(uint(couponID), Init.DB)
	if errrrr != nil {
		c.JSON(400, gin.H{"Error": "Fetching coupon by id error"})
		return
	}

	currentTime := time.Now()
	if coupon.ExpiresAt.Before(currentTime) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon expired"})
		return
	}

	if amount < coupon.MinValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum amount required"})
		return
	}

	if amount > coupon.MaxValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount above max amount"})
		return
	}

	result, _ := fetchUsedCoupon(coupon.ID, user.ID, Init.DB)
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon already used"})
		return
	}

	updatedTotal := amount - coupon.Discount
	err = setRedis("Amount", updatedTotal, 1*time.Hour)
	err = setRedis("couponID", couponID, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "coupon applied",
		"couponDiscount": coupon.Discount,
		"current total":  updatedTotal,
	})
}

// ViewWallet retrieves user's wallet information.
func ViewWallet(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	user, err := fetchUser(username, Init.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching user"})
		return
	}

	wallet, err := fetchWallet(user.ID, Init.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching wallet"})
		return
	}

	err = setRedis("WalletId", wallet.ID, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": wallet})
}

// ApplyWallet applies user's wallet balance.
func ApplyWallet(c *gin.Context) {
	amountStr, err := getRedis("Amount")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'amount' from Redis client"})
		return
	}
	amountt, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
		return
	}
	walletIDStr, err := getRedis("WalletId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'walletId' from Redis client"})
		return
	}
	walletID, err := strconv.Atoi(walletIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
		return
	}

	wallet, err := fetchWalletByID(uint(walletID), Init.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching wallet"})
		return
	}

	amount := float64(amountt)
	if wallet.Balance <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no balance available"})
		return
	}

	var balance float64
	if amount > wallet.Balance {
		balance = 0
	} else {
		balance = amount - wallet.Balance
	}
	amount = amount - wallet.Balance
	wallet.Balance = balance
	errr := saveWallet(Init.DB)
	if errr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while updating wallet"})
		return
	}


	transaction.Date = time.Now()
	transaction.Details = "Booked room in"
	transaction.Amount = amount
	transaction.UserID = wallet.UserID

	errrr := createwallet(Init.DB)
	if errrr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction adding error"})
		return
	}

	err = setRedis("Amount", amount, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "wallet applied",
		"wallet balance": wallet.Balance,
		"amount":         amount,
	})
}
