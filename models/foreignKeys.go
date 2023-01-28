package models

import "time"

type Dealer struct {
	ID            int64      `json:"id" gorm:"column:id;primary_key;autoIncrement" form:"id"`
	PartnerID     int64      `json:"partner_id" gorm:"column:partner_id"`
	Partner       Partner    `gorm:"foreignkey:PartnerID" json:"-"`
	Limit         float64    `json:"limit" gorm:"column:limit"`
	Distributed   float64    `json:"distributed" gorm:"column:distributed"`
	INN           string     `json:"inn" gorm:"column:inn"`
	Name          string     `json:"name" gorm:"column:name"`
	Active        bool       `json:"active" gorm:"column:active"`
	Address       string     `json:"address" gorm:"column:address"`
	RegionID      int64      `json:"region_id" gorm:"column:region_id"`
	CityID        int64      `json:"city_id" gorm:"column:city_id"`
	CreatedAt     *time.Time `gorm:"autoCreateTime"`
	UpdatedAt     *time.Time `gorm:"autoUpdateTime"`
	CreatedBy     int64      `json:"created_by" gorm:"column:created_by"`
	UpdatedBy     int64      `json:"updated_by" gorm:"column:updated_by"`
	Director      string     `json:"director" gorm:"column:director"`
	Phone         string     `json:"phone" gorm:"column:phone"`
	CuratorID     int64      `json:"curator_id" gorm:"column:curator_id"`
	DealerCftId   string     `json:"dealer_cft_id" gorm:"column:dealer_cft_id"`
	SuperAgent    bool       `json:"super_agent" gorm:"column:super_agent"`
	ActivityBasis string     `json:"activity_basis" gorm:"column:activity_basis"`
	OrzuCashOut   bool       `json:"orzu_cash_out" gorm:"column:orzu_cash_out"`
}

type Terminal struct {
	ID            int64      `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	DealerID      int64      `json:"dealerID" gorm:"column:dealer_id"`
	Dealer        Dealer     `gorm:"foreignkey:DealerID" json:"-"`
	Active        bool       `json:"active" gorm:"column:active"`
	Address       string     `json:"address" gorm:"column:address"`
	RegionID      int64      `json:"region_id" gorm:"column:region_id"`
	CityID        int64      `json:"city_id" gorm:"column:city_id"`
	Balance       float64    `json:"balance" gorm:"column:balance"`
	Limit         float64    `json:"limit" gorm:"column:limit"`
	CashOut       bool       `json:"cash_out" gorm:"column:cash_out;default:false"`
	CftTerminalID string     `json:"cft_terminal_id" gorm:"column:cft_terminal_id"`
	CreatedAt     *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy     int64      `json:"created_by" gorm:"column:created_by"`
	UpdatedBy     int64      `json:"updated_by" gorm:"column:updated_by"`
	SuperTerminal bool       `json:"super_terminal" gorm:"column:super_terminal"`
	OrzuCashOut   bool       `json:"orzu_cash_out" gorm:"column:orzu_cash_out"`
}

type Partner struct {
	ID   int64  `gorm:"column:id;primary_key;autoIncrement"`
	Name string `gorm:"column:name"`
}

type TUser struct {
	ID           int64      `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	DealerID     int64      `json:"dealerID" gorm:"column:dealer_id" validate:"gte=1"`
	TerminalID   int64      `json:"terminalID" gorm:"column:terminal_id" validate:"gte=1"`
	RoleID       int64      `json:"roleID" gorm:"column:role_id" validate:"gte=1"`
	BranchID     int64      `json:"branch_id" gorm:"column:branch_id" validate:"gte=1"`
	RegionID     int64      `json:"region_id" gorm:"column:region_id"`
	CityID       int64      `json:"city_id" gorm:"column:city_id"`
	FIO          string     `json:"fio" gorm:"column:fio"`
	Phone        string     `json:"phone" gorm:"column:phone"`
	Login        string     `json:"login" gorm:"column:login" validate:"gt=0"`
	Password     string     `json:"password" gorm:"column:password" validate:"gt=0"`
	Salt         string     `json:"salt" gorm:"column:salt"`
	RefreshToken string     `json:"refresh_token" gorm:"column:refresh_token"`
	Active       bool       `json:"active" gorm:"column:active"`
	CreatedAt    *time.Time `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime"`
	LoginAt      *time.Time `gorm:"column:login_at"`
}
