package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// HotelCategory model
type HotelCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}

// HotelAmenities model
type HotelAmenities struct {
	FacilityID     uint   `json:"facility_id" gorm:"primaryKey;autoIncrement"`
	HotelAmenities string `json:"amenities" gorm:"not null"`
}

// JSONB type for handling JSON data in the database
type JSONB []interface{}

// Value used to retrive value
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan helps to scan values
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// Hotels model
type Hotels struct {
	gorm.Model
	Name            string        `json:"name" validate:"required"`
	Title           string        `json:"title" validate:"required"`
	Description     string        `json:"description" validate:"required"`
	StartingPrice   float64       `json:"starting_price" validate:"required"`
	City            string        `json:"city" validate:"required"`
	Pincode         string        `json:"pincode" validate:"required"`
	Address         string        `json:"address" validate:"required"`
	Images          string        `json:"images" validate:"required"`
	TypesOfRoom     int           `json:"types_of_room" validate:"required"`
	Facility        JSONB         `gorm:"type:jsonb" json:"facilities"`
	Revenue         float64       `json:"revenue" gorm:"default=0"`
	IsAvailable     bool          `json:"is_available" gorm:"default=false"`
	IsBlock         bool          `json:"is_block"`
	AdminApproval   bool          `json:"admin_approval" gorm:"default=false"`
	HotelCategoryID uint          `json:"category_id" gorm:"not null"`
	HotelCategory   HotelCategory `gorm:"ForeignKey:HotelCategoryID"`
	OwnerUsername   string
}

func (hotel *Hotels) FetchHotels(cityORhotelname string, isAvailab, isBlock, adminApproval bool, skip, limit int, db *gorm.DB) ([]Hotels, error) {
	var hotels []Hotels
	if err := db.Preload("HotelCategory").Offset(skip).Limit(limit).Where("(city ILIKE ? OR name ILIKE ?) AND is_available = ? AND is_block = ? AND admin_approval = ?", "%"+cityORhotelname+"%", "%"+cityORhotelname+"%", isAvailab, isBlock, adminApproval).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

func (hotel *Hotels) FetchHotelById(hotelId uint, db *gorm.DB) (*Hotels, error) {
	return hotel, nil
}

func (hotel *Hotels) FetchAvailableHotels(loc string, db *gorm.DB) ([]Hotels, error) {
	var hotels []Hotels
	if err := db.Where("city = ?", loc).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}
