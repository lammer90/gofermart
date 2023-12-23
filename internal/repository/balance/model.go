package balance

type BalanceRepository interface {
	CreateBalance(login string) error
	AddBonus(login string, sumToAdd float32) error
	WithdrawBonus(login string, sumToMinus float32) error
	FindByUser(login string) (*Balance, error)
}

type Balance struct {
	Login       string
	BonusSum    float32
	WithdrawSum float32
}
