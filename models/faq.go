package models

import (
	"time"

	"gorm.io/gorm"
)

type Faq struct {
	ID         int64           `json:"id" gorm:"column:id;primary_key;autoIncrement" example:"1"`
	Key        string          `json:"key" gorm:"column:key" example:"1"`
	Question   string          `json:"question" gorm:"column:question" example:"1"`
	Response   string          `json:"response" gorm:"column:response" example:"1"`
	CategoryID int64           `json:"category_id" gorm:"column:category_id" example:"1"`
	CreatedAt  *time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  *time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  *gorm.DeletedAt `json:"-" gorm:"index"`
	Category   FaqCategory     `json:"-" gorm:"foreignkey:CategoryID"`
}

func (Faq) TableName() string {
	return "faqs"
}

type FaqCategory struct {
	ID        int64           `json:"id" gorm:"column:id;primary_key;autoIncrement" example:"1"`
	Name      string          `json:"name" gorm:"column:name" example:"1"`
	Product   string          `json:"product" gorm:"column:product" example:"1"`
	Project   string          `json:"project" gorm:"column:project" example:"1"`
	CreatedAt *time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt *time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"index"`
}

func (FaqCategory) TableName() string {
	return "faq_categories"
}

type FaqFilters struct {
	ID         *int64  `form:"id" example:"1"`
	Key        *string `form:"key" example:"1"`
	Search     *string `form:"search" example:"1"`
	Question   *string `form:"question" example:"1"`
	CategoryID *int64  `form:"category_id" example:"1"`
	Product    *string `form:"product" example:"1"`
	Project    *string `form:"project" example:"1"`
	Page       int64   `form:"page" example:"1"`
	PageLimit  int64   `form:"page_limit" example:"1"`
}

type LandingData struct {
	FullName  string `gorm:"column:full_name" json:"full_name"` //фио (не обязательное поле)
	Phone     string `gorm:"column:phone" json:"phone"`         //номер тел (обязательное поле)
	Entity    string `gorm:"column:entity" json:"entity"`       // организация/ИП (не обязательное поле)
	City      string `gorm:"column:city" json:"city"`           //город/район (обязательное поле)
	Delivered bool   `gorm:"column:delivered" json:"delivered"` //это для системы(бэк) флажок для бота, после отправки сделать true!
	IP        string `gorm:"column:ip" json:"ip"`               //это для системы(бэк)
}

func (LandingData) TableName() string {
	return "landing_data"
}
