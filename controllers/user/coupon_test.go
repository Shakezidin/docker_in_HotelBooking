package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCalculateAmountForDays(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
	}{{
		name:        "CalculateAmountForDays Success",
		route:       "/user/home/room/book?id=1",
		errorResult: nil,
	}, {
		name:        "CalculateAmountForDays params error",
		route:       "/user/home/room/book?id=",
		errorResult: map[string]string{"error": "roomid query parameter is missing"},
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			setRedis = func(key string, value any, expirationTime time.Duration) error {
				return nil
			}
			getRedis = func(key string) (string, error) {
				return "2006-01-02", nil
			}
			fetchRoomById = func(roomID uint, db *gorm.DB) (*models.Rooms, error) {
				return &rooms, nil
			}

			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.route, nil)
			if err != nil {
				require.NoError(t, err)
			}
			engine.ServeHTTP(w, req)
			if tc.errorResult != nil {
				errValue, _ := json.Marshal(tc.errorResult)
				require.Equal(t, w.Body.String(), string(errValue))
			} else {
				require.Equal(t, w.Code, 200)
			}
		})
	}
}

func TestViewNonBlockedCoupons(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
	}{{
		name:        "View Coupon Success",
		route:       "/user/coupons/view",
		errorResult: nil,
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchAllCoupon = func(db *gorm.DB) ([]models.Coupon, error) {
				var coupons []models.Coupon
				return coupons, nil
			}
			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.route, nil)
			if err != nil {
				require.NoError(t, err)
			}
			engine.ServeHTTP(w, req)
			require.Equal(t, w.Code, 200)
		})
	}
}

func TestApplyCoupon(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
	}{{
		name:        "error",
		route:       "/user/coupon/apply?id=",
		errorResult: map[string]string{"Error": "query id error"},
	}, {
		name:        "applayCoupon Success",
		route:       "/user/coupon/apply?id=1",
		errorResult: nil,
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			getRedis = func(key string) (string, error) {
				return "500", nil
			}

			fetchUser = func(val string, db *gorm.DB) (*models.User, error) {
				return &user, nil
			}

			fetchCouponById = func(couponId uint, db *gorm.DB) (*models.Coupon, error) {
				dateStr := "2023-12-12"
				parsedTime, _ := time.Parse("2006-01-02", dateStr)
				var coupon models.Coupon
				coupon.ID = 1
				coupon.CouponCode = "sample code"
				coupon.Discount = 200
				coupon.ExpiresAt = parsedTime
				coupon.MaxValue = 10000
				coupon.MinValue = 10
				return &coupon, nil
			}

			fetchUsedCoupon = func(couponId, userId uint, db *gorm.DB) (*models.UsedCoupon, error) {
				return nil, nil
			}
			setRedis = func(key string, value any, expirationTime time.Duration) error {
				return nil
			}
			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.route, nil)
			if err != nil {
				require.NoError(t, err)
			}
			req.Header.Set("Authorization", authToken)
			engine.ServeHTTP(w, req)
			if tc.errorResult != nil {
				errReslt, err := json.Marshal(tc.errorResult)
				if err != nil {
					require.NoError(t, err)
				}
				require.Equal(t, w.Body.String(), string(errReslt))
			} else {
				require.Equal(t, w.Code, 200)
			}
		})
	}
}

func TestViewWallet(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
	}{{
		name:        "ViewWalletSuccess",
		route:       "/user/wallet",
		errorResult: nil,
	}}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			fetchUser = func(val string, db *gorm.DB) (*models.User, error) {
				return &user, nil
			}
			fetchWallet = func(userId uint, db *gorm.DB) (*models.Wallet, error) {
				return &walletref, nil
			}
			gin.SetMode(gin.TestMode)
			engine := gin.Default()
			RegisterUserRoutes(engine)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.route, nil)
			if err != nil {
				require.NoError(t, err)
			}
			req.Header.Set("Authorization", authToken)
			engine.ServeHTTP(w, req)
			if tc.errorResult != nil {
				errReslt, err := json.Marshal(tc.errorResult)
				if err != nil {
					require.NoError(t, err)
				}
				require.Equal(t, w.Body.String(), string(errReslt))
			} else {
				require.Equal(t, w.Code, 200)
			}
		})
	}
}

func TestApplyWallet(t *testing.T) {
	test := []struct {
		name        string
		route       string
		errorResult map[string]string
	}{{
		name:        "Applaywallet Success",
		route:       "/user/wallet/apply",
		errorResult: nil,
	}}
	for _, tc := range test {
		getRedis = func(key string) (string, error) {
			return "500", nil
		}
		fetchWalletByID = func(walletId uint, db *gorm.DB) (*models.Wallet, error) {
			var wallet models.Wallet
			wallet.Balance = 500
			return &wallet, nil
		}
		saveWallet = func(db *gorm.DB) error {
			return nil
		}
		setRedis = func(key string, value any, expirationTime time.Duration) error {
			return nil
		}
		createwallet = func(db *gorm.DB) error {
			return nil
		}
		gin.SetMode(gin.TestMode)
		engine := gin.Default()
		RegisterUserRoutes(engine)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, tc.route, nil)
		if err != nil {
			require.NoError(t, err)
		}
		req.Header.Set("Authorization", authToken)
		engine.ServeHTTP(w, req)
		if tc.errorResult != nil {
			errReslt, err := json.Marshal(tc.errorResult)
			if err != nil {
				require.NoError(t, err)
			}
			require.Equal(t, w.Body.String(), string(errReslt))
		} else {
			require.Equal(t, w.Code, 200)
		}
	}
}
