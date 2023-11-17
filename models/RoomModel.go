package models

import (
	"time"

	"gorm.io/gorm"
)

// Cancellation Model
type Cancellation struct {
	gorm.Model
	CancellationPolicy     string `json:"cancellation_policy" gorm:"not null"`
	RefundAmountPercentage int    `json:"refund_amount_percentage" gorm:"not null"`
}

// AvailableRoom Model
type AvailableRoom struct {
	gorm.Model
	RoomID    uint `json:"room_id"`
	BookingID uint
	CheckIn   time.Time `json:"check_in" time_format:"2006-01-02"`
	CheckOut  time.Time `json:"check_out" time_format:"2006-01-02"`
}

// RoomFacilities Model
type RoomFacilities struct {
	gorm.Model
	RoomAmenities string `json:"room_amenities" gorm:"not null" validate:"required"`
}

// RoomCategory Model
type RoomCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

func (roomCategory *RoomCategory) FetchRoomCategory(db *gorm.DB) ([]RoomCategory, error) {
	var categories []RoomCategory
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Rooms Model
type Rooms struct {
	gorm.Model
	Description    string       `json:"description" validate:"required" gorm:"not null"`
	Price          float64      `json:"price" gorm:"not null" validate:"required"`
	Adults         int          `json:"adults" gorm:"not null" validate:"required"`
	Children       int          `json:"children" gorm:"not null" validate:"required"`
	Bed            string       `json:"bed" gorm:"not null" validate:"required"`
	Images         string       `json:"images" validate:"required"`
	CancellationID uint         `json:"cancellation_id" gorm:"not null"`
	Cancellation   Cancellation `gorm:"ForeignKey:CancellationID"`
	Facility       JSONB        `gorm:"type:jsonb" json:"facilities"`
	RoomNo         int
	IsAvailable    bool         `json:"is_available" validate:"required"`
	IsBlocked      bool         `json:"is_blocked"`
	DiscountPrice  float64      `json:"discount_price"`
	Discount       float64      `json:"discount"`
	AdminApproval  bool         `json:"admin_approval" gorm:"default=false"`
	HotelsID       uint         `json:"hotel_id" gorm:"not null"`
	Hotels         Hotels       `gorm:"ForeignKey:HotelsID"`
	OwnerUsername  string       `json:"owner_username"`
	RoomCategoryID uint         `json:"category_id" gorm:"not null"`
	RoomCategory   RoomCategory `json:"category" gorm:"ForeignKey:RoomCategoryID"`
}

func (rooms *Rooms) FetchinRooms(hotelId uint, db *gorm.DB) ([]Rooms, error) {
	var roomss []Rooms
	if err := db.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("hotels_id = ? AND is_available = ? AND is_blocked = ? AND admin_approval = ?", hotelId, true, false, true).Find(&roomss).Error; err != nil {
		return nil, err
	}
	return roomss, nil
}

func (rooms *Rooms) FetchingAvailableRooms(hotelId uint, fromDate, toDate time.Time, NumberOfAdults, NumberOfChildren int, db *gorm.DB) ([]Rooms, error) {
	var availableRooms []Rooms
	if err := db.Where("hotels_id = ? AND (available_rooms.room_id IS NULL OR ? > available_rooms.check_out OR ? < available_rooms.check_in)",
		hotelId, fromDate, toDate).Where("adults >= ? AND children >= ? AND is_blocked = ? AND admin_approval = ?", NumberOfAdults, NumberOfChildren, false, true).
		Joins("LEFT JOIN available_rooms ON rooms.id = available_rooms.room_id").
		Joins("LEFT JOIN room_categories ON rooms.room_category_id = room_categories.id").
		Preload("RoomCategory").
		Find(&availableRooms).Error; err != nil {
		return nil, err
	}
	return availableRooms, nil
}
func (rooms *Rooms) FetchAllRooms(skip, limit int, db *gorm.DB) ([]Rooms, error) {
	var roomss []Rooms
	if err := db.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Offset(skip).Limit(limit).Where("is_available = ? AND is_blocked = ? AND admin_approval = ?", true, false, true).Find(&roomss).Error; err != nil {
		return nil, err
	}
	return roomss, nil
}

func (rooms *Rooms) FetchRoomById(roomID uint, db *gorm.DB) (*Rooms, error) {
	if err := db.Preload("Hotels").Preload("Cancellation").Preload("RoomCategory").
		First(&rooms, uint(roomID)).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}
