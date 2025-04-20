package schema

import (
	"time"

	"gorm.io/gorm"

	"ar-app-api/model"
)

type User struct {
	gorm.Model
	Username    string       `gorm:"uniqueIndex;size:64;NOT NULL"`  // 用户名，唯一索引
	Password    string       `gorm:"size:128;NOT NULL"`             // 密码，通常加密存储
	Email       string       `gorm:"uniqueIndex;size:128;NOT NULL"` // 邮箱，唯一索引
	PhoneNumber string       `gorm:"size:20;NOT NULL"`              // 手机号码
	FirstName   string       `gorm:"size:64"`                       // 名
	LastName    string       `gorm:"size:64"`                       // 姓
	Gender      string       `gorm:"size:10"`                       // 性别
	DateOfBirth time.Time    // 出生日期
	Address     string       `gorm:"size:256"` // 地址
	City        string       `gorm:"size:64"`  // 城市
	State       string       `gorm:"size:64"`  // 州/省
	Country     string       `gorm:"size:64"`  // 国家
	PostalCode  string       `gorm:"size:20"`  // 邮政编码
	Status      string       `gorm:"size:20"`  // 账户状态 (例如：active, inactive, banned)
	LastLogin   time.Time    // 最后登录时间
	ChatGroups  []*ChatGroup `gorm:"many2many:group_members;"`
}

func (user *User) ToType() *model.User {
	groups := make([]*model.TinyGroup, 0, len(user.ChatGroups))
	for _, g := range user.ChatGroups {
		groups = append(groups, &model.TinyGroup{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			CreatorID:   g.CreatorID,
			CreatedAt:   g.CreatedAt,
			IsPrivate:   g.IsPrivate,
			UpdatedAt:   g.UpdatedAt,
		})
	}
	return &model.User{
		ID:          user.ID,
		Username:    user.Username,
		Password:    user.Password,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
		Address:     user.Address,
		City:        user.City,
		State:       user.State,
		Country:     user.Country,
		PostalCode:  user.PostalCode,
		Status:      user.Status,
		Groups:      groups,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLogin:   user.LastLogin,
	}
}
