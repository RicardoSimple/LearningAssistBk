package model

import (
	"time"
)

type User struct {
	ID           uint         `json:"ID,omitempty"`
	Username     string       `json:"username,omitempty"`
	Password     string       `json:"-"`
	Email        string       `json:"email,omitempty"`
	PhoneNumber  string       `json:"phoneNumber,omitempty"`
	FirstName    string       `json:"firstName,omitempty"`
	LastName     string       `json:"lastName,omitempty"`
	Gender       string       `json:"gender,omitempty"`
	DateOfBirth  time.Time    `json:"dateOfBirth"`
	Address      string       `json:"address,omitempty"`
	City         string       `json:"city,omitempty"`
	State        string       `json:"state,omitempty"`
	Country      string       `json:"country,omitempty"`
	PostalCode   string       `json:"postalCode,omitempty"`
	Status       string       `json:"status,omitempty"`
	Groups       []*TinyGroup `json:"groups,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	LastLogin    time.Time    `json:"lastLogin"`
	UserType     string       `json:"userType,omitempty"`
	ClassStage   string       `json:"class_stage,omitempty"`
	Name         string       `json:"name,omitempty"`
	AddClassTime time.Time    `json:"addClassTime,omitempty"`
	ClassId      uint         `json:"classId,omitempty"`
}

// TinyGroup 基础群聊信息
type TinyGroup struct {
	ID          uint
	Name        string
	Description string
	CreatorID   uint      // 群聊创建者的用户 ID
	CreatedAt   time.Time // 群聊创建时间
	IsPrivate   bool      // 是否为私有群聊
	UpdatedAt   time.Time
}
