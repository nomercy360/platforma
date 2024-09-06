package admin

import (
	"rednit/config"
	"rednit/db"
)

type storage interface {
	GetUser(db.UserQuery) (*db.User, error)
	CreateUser(string, string, *string) (*db.User, error)
	ListCustomers() ([]db.Customer, error)
	ListDiscounts() ([]db.Discount, error)
	ListOrders() ([]db.Order, error)
	ListProducts(params db.ListProductsQuery) ([]db.Product, error)
	ListUsers() ([]db.User, error)
}

type Admin struct {
	s   storage
	cfg config.Default
}

func New(s storage, cfg config.Default) Admin {
	return Admin{s: s, cfg: cfg}
}
