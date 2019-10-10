package account

import (
	"github.com/IntouchOpec/base-go-echo/model"
)

type Account interface {
	Get(id string) interface{}
	Set(Account interface{}) interface{}
}

type account struct {
	model.Account
}

func (a *account) Get(id string) *account {
	model.DB().Find(&a, id)
	return a
}
