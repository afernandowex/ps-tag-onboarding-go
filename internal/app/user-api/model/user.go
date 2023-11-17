package model

type User struct {
	ID        int32  `gorm:"primary_key" json:"id"`
	FirstName string `gorm:"not null" json:"firstName"`
	LastName  string `gorm:"not null" json:"lastName"`
	Email     string `gorm:"not null" json:"email"`
	Age       int8   `gorm:"not null" json:"age"`
}
