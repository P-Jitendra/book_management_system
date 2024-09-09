package models

import "time"

type UserInfo struct {
	UserId    string `gorm:"primary_key;not null"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"`
	ContactNo int64  `gorm:"uniqueIndex;not null"`
	Category  string `gorm:"default:user;not null"`
}

type BooksInfo struct {
	BookId       int    `gorm:"primary_key;autoIncrement;not null"`
	Title        string `gorm:"not null"`
	Author       string `gorm:"type:varchar(255);not null"`
	Publications string `gorm:"not null"`
	Genre        string `gorm:"default:novels;not null"`
}

type LendBooksInfo struct {
	BookId            int       `gorm:"primaryKey"`
	UserId            string    `gorm:"primaryKey"`
	Status            string    `gorm:"not null"`
	BorrowedTimestamp time.Time `gorm:"primaryKey"`
	ReturnedTimestamp time.Time `gorm:"default:NULL"`
}
