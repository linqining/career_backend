package model

type Company struct {
	ID            int64  `gorm:"column:id;primary_key" json:"id"`             //
	Name          string `gorm:"column:name" json:"name"`                     //
	WalletAddress string `gorm:"column:wallet_address" json:"wallet_address"` //
}

// TableName sets the insert table name for this struct type
func (c *Company) TableName() string {
	return "company"
}
