package service

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/repository"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepository repository.WalletRepositorier
	Database         *gorm.DB
}

func (d WalletService) GetWallet(ctx context.Context, address string) (*model.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (d WalletService) Transfer(ctx context.Context, fromAddress string, toAddress string, amount string) (string, error) {
	//TODO implement me
	panic("implement me")
}
