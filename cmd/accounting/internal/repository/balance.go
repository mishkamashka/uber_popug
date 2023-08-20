package repository

import (
	"time"
	"uber-popug/pkg/types"
)

type Balance struct {
	UserID    string    `gorm:"owner_id;primarykey"`
	Amount    int       `gorm:"amount"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func BalanceToRepoType(u *types.Balance) *Balance {
	return &Balance{
		UserID:    u.UserID,
		Amount:    u.Amount,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func RepoTypeToBalance(u *Balance) *types.Balance {
	return &types.Balance{
		UserID:    u.UserID,
		Amount:    u.Amount,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func RepoTypesToBalances(u []*Balance) []*types.Balance {
	res := make([]*types.Balance, 0, len(u))

	for _, user := range u {
		resUser := RepoTypeToBalance(user)

		res = append(res, resUser)
	}

	return res
}
