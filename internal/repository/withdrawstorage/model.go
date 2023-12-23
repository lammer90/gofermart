package withdrawstorage

import (
	"github.com/lammer90/gofermart/internal/dto/withdraw"
)

type WithdrawRepository interface {
	Save(withdraw *withdraw.Withdraw) error
	FindByUser(login string) ([]withdraw.Withdraw, error)
	FindSumByUser(login string) (float32, error)
}
