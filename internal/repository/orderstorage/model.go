package orderstorage

import (
	"github.com/lammer90/gofermart/internal/dto/order"
)

type OrderRepository interface {
	Save(order *order.Order) error
	FindByUser(login string) ([]order.Order, error)
	FindByNumber(number string) (*order.Order, error)
	FindNumbersToProcess() ([]string, error)
	Update(order *order.Order) error
}
