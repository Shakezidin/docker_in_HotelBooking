package models

import "gorm.io/gorm"

// Contact Model
type Contact struct {
	gorm.Model
	Message string `json:"message"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"ForeignKey:UserID"`
}

func (Contact *Contact) CreateContact(db *gorm.DB) error {
	result := db.Create(&Contact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Createcontact(contact *Contact,db *gorm.DB)error{
	return contact.CreateContact(db)
}

type Message struct {
	Message string `json:"message"`
}