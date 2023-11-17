package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Coupon Model
type Coupon struct {
	gorm.Model
	CouponCode string    `gorm:"not null"`
	Discount   int       `gorm:"not null"`
	MinValue   int       `gorm:"not null"`
	MaxValue   int       `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	IsBlock    bool      `gomr:"default:false"`
}

// UsedCoupon Model
type UsedCoupon struct {
	gorm.Model
	UserID   uint
	CouponID uint
}

func (coupon *Coupon) FetchAllCoupon(db *gorm.DB) ([]Coupon, error) {
	var coupons []Coupon
	if err := db.Where("is_block = ?", false).Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

func (coupon *Coupon) FetchCouponById(couponId uint, db *gorm.DB) (*Coupon, error) {
	if err := db.Where("is_block = ? And id = ?", false, couponId).First(&coupon).Error; err != nil {
		return nil, err
	}
	return coupon, nil
}

func (usedCoupon *UsedCoupon) FetchUsedCoupon(couponId, userId uint, db *gorm.DB) (*UsedCoupon, error) {
	user := db.Where("coupon_id = ? AND user_id = ?", couponId, userId).First(usedCoupon)
	if errors.Is(user.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No record found, no error
	}
	if user.Error != nil {
		return nil, user.Error // Other errors
	}
	return usedCoupon, nil
}
