package service

import (
	"context"
	"errors"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/repository"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepository repository.WalletRepositorier
	Database         *gorm.DB
}

func (d *WalletService) GetWallet(ctx context.Context, address string) (*model.Wallet, error) {
	return d.WalletRepository.GetWalletByAddress(ctx, d.Database, address)
}

func (d *WalletService) Transfer(ctx context.Context, fromAddress string, toAddress string, amount int) (int, error) {
	//TODO: avoid deadlocks (smaller/bigger address first)
	var newBalance int

	err := d.Database.Transaction(func(tx *gorm.DB) error {
		fromWallet, intErr := d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, fromAddress)
		if intErr != nil {
			return intErr // rollback on any error
		}

		if fromWallet.Tokens < amount {
			return errors.New("insufficient balance")
		}

		toWallet, intErr := d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, toAddress)
		if intErr != nil {
			return intErr
		}

		newFromWalletBalance := fromWallet.Tokens - amount
		newToWalletBalance := toWallet.Tokens + amount

		intErr = d.WalletRepository.UpdateWalletTokensByAddress(ctx, tx, fromAddress, newFromWalletBalance)
		if intErr != nil {
			return intErr
		}

		intErr = d.WalletRepository.UpdateWalletTokensByAddress(ctx, tx, toAddress, newToWalletBalance)
		if intErr != nil {
			return intErr
		}

		newBalance = newFromWalletBalance
		return nil
	})
	if err != nil {
		return -1, err
	}

	return newBalance, nil
}
