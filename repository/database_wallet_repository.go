package repository

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"gorm.io/gorm"
)

type DatabaseWalletRepository struct {
}

func (d DatabaseWalletRepository) GetWalletByAddress(ctx context.Context, tx *gorm.DB, address string) (*model.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (d DatabaseWalletRepository) GetWalletByAddressForUpdate(ctx context.Context, tx *gorm.DB, address string) (*model.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (d DatabaseWalletRepository) UpdateWalletByAddress(ctx context.Context, tx *gorm.DB, address string, wallet *model.Wallet) error {
	//TODO implement me
	panic("implement me")
}
