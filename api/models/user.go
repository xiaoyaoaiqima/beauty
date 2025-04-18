package models

import (
	"time"
)

// User 用户模型
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
    Password  string    `json:"password,omitempty" gorm:"type:varchar(255);not null"`
    Email     string    `json:"email" gorm:"type:varchar(100);uniqueIndex"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}