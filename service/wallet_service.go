package service

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/repository"
)

type WalletService struct {
	WalletRepository repository.WalletRepositorier
}

func (d WalletService) GetWalletByAddress(ctx context.Context, address string) (*model.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (d WalletService) Transfer(ctx context.Context, fromAddress string, toAddress string, amount string) (string, error) {
	//TODO implement me
	panic("implement me")
}
