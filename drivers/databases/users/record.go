package users

import (
	"fgd-alterra-29/business/users"
)

type Users struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}

func (user *Users) ToDomain() users.Domain {
	return users.Domain{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
	}
}

func ToListDomain(u []Users) []users.Domain {
	var Domains []users.Domain

	for _, val := range u {
		Domains = append(Domains, val.ToDomain())
	}
	return Domains
}
