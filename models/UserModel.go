package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	UserName     string    `json:"username" gorm:"not null;unique" validate:"required"`
	Name         string    `json:"name" gorm:"not null" validate:"required"`
	Email        string    `json:"email" gorm:"not null;unique" validate:"required"`
	Phone        string    `json:"phone" gorm:"not null;unique" validate:"required"`
	Password     string    `json:"password" gorm:"not null" validate:"required"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false"`
	Wallet       Wallet    `json:"wallet"`
	ReferralCode string    `json:"referral_code"`
	JoinedAt     time.Time `json:"joined_at" gorm:"default:now()"`
}

// func (user *User) CreateUser(db *gorm.DB) error {
// 	if err := db.Create(&user); err != nil {
// 		return err.Error
// 	}
// 	return nil
// }

func (user *User) CreateUser(userr *User, db *gorm.DB) error {
	result := db.Create(&userr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) FetchUser(val string, db *gorm.DB) (*User, error) {
	if err := db.Where("user_name = ?", val).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (user *User) UpdateUser(db *gorm.DB) error {
	save := db.Save(&user)
	if save.Error != nil {
		return save.Error
	}
	return nil
}

func UpdateUSer(user *User, db *gorm.DB) error {
	return user.UpdateUser(db)
}

// func Userfetch(db *gorm.DB,user *User)(*User,error){
// 	return user.FetchUser(db)
// }

// Wallet represents a user's wallet balance.
type Wallet struct {
	gorm.Model
	Balance float64
	UserID  uint `gorm:"unique"`
}

// Transaction represents a financial transaction.
type Transaction struct {
	gorm.Model
	Date    time.Time
	Details string
	Amount  float64
	UserID  uint
}

func (transaction *Transaction) Create(db *gorm.DB) error {
	if err := db.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}

type OtpCredentials struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func (wallet *Wallet) FetchWallet(userId uint, db *gorm.DB) (*Wallet, error) {
	if err := db.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *Wallet) FetchWalletById(walletId uint, db *gorm.DB) (*Wallet, error) {
	if err := db.Where("id = ? = ?", walletId).First(&wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *Wallet) SaveWallet(db *gorm.DB) error {
	if err := db.Save(&wallet).Error; err != nil {
		return err
	}
	return nil
}
