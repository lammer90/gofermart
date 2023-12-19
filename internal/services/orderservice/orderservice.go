package orderservice

import (
	"errors"
	"github.com/EClaesson/go-luhn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lammer90/gofermart/internal/dto/order"
	"github.com/lammer90/gofermart/internal/repository/orderstorage"
	"time"
)

type orderServiceImpl struct {
	repository orderstorage.OrderRepository
}

func New(repository orderstorage.OrderRepository) OrderService {
	return &orderServiceImpl{repository: repository}
}

func (o orderServiceImpl) Save(number, login string) error {
	valid, err := luhn.IsValid(number)
	if err != nil || !valid {
		return NotValidLuhnSum
	}

	existedOrder, err := o.repository.FindByNumber(number)
	if err != nil {
		return err
	}
	if existedOrder != nil {
		if login == existedOrder.Login {
			return OrderNumberAlreadyExistThisUser
		}
		return OrderNumberAlreadyExistAnotherUser
	}

	newOrder := &order.Order{Login: login, Number: number, Status: order.NEW, Accrual: 0.00, UploadedAt: time.Now()}
	err = o.repository.Save(newOrder)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return OrderNumberAlreadyExistAnotherUser
		}
		return err
	}
	return nil
}

func (o orderServiceImpl) FindAll(login string) ([]order.OrderResponse, error) {
	response := make([]order.OrderResponse, 0)

	orders, err := o.repository.FindByUser(login)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return nil, nil
	}
	for _, ord := range orders {
		resp := order.OrderResponse{
			ord.Number,
			order.Statuses[ord.Status-1],
			ord.Accrual,
			ord.UploadedAt.Format("2006-01-02T15:04:05-07:00")}
		response = append(response, resp)
	}
	return response, nil
}
