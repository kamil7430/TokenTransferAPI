package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseWalletRepository struct {
}

func (d *DatabaseWalletRepository) GetWalletByAddress(ctx context.Context, tx *gorm.DB, address string) (*model.Wallet, error) {
	wallet, err := gorm.G[model.Wallet](tx).Where("Address = ?", address).First(ctx)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (d *DatabaseWalletRepository) GetWalletByAddressForUpdate(ctx context.Context, tx *gorm.DB, address string) (*model.Wallet, error) {
	wallet, err := gorm.G[model.Wallet](tx, clause.Locking{Strength: "UPDATE"}).Where("Address = ?", address).First(ctx)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (d *DatabaseWalletRepository) UpdateWalletTokensByAddress(ctx context.Context, tx *gorm.DB, address string, tokens string) error {
	rows, err := gorm.G[model.Wallet](tx).Where("Address = ?", address).Update(ctx, "Tokens", tokens)
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New(fmt.Sprintf("affected %d rows, expected 1", rows))
	}
	return nil
}
