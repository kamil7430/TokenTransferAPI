package repository

import (
	"context"
	"errors"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
)

type InMemoryWalletRepository struct {
	wallets map[string]*model.Wallet
}

func (i InMemoryWalletRepository) GetWalletByAddress(ctx context.Context, address string) (*model.Wallet, error) {
	val, ok := i.wallets[address]
	if !ok {
		return nil, errors.New("this wallet does not exist")
	}
	return val, nil
}
