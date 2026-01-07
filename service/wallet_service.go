package service

import (
	"context"
	"errors"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/helper/address_helper"
	"github.com/kamil7430/TokenTransferAPI/repository"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepository repository.WalletRepositorier
	Database         *gorm.DB
}

func (d *WalletService) GetWallet(ctx context.Context, address string) (*model.Wallet, error) {
	err := address_helper.CheckAddress(address)
	if err != nil {
		return nil, err
	}
	return d.WalletRepository.GetWalletByAddress(ctx, d.Database, address)
}

func (d *WalletService) Transfer(ctx context.Context, fromAddress string, toAddress string, amount int) (int, error) {
	if amount <= 0 {
		return -1, errors.New("amount must be greater than zero")
	}
	if fromAddress == toAddress {
		return -1, errors.New("from and to addresses cannot be equal")
	}

	err := address_helper.CheckAddress(fromAddress)
	if err != nil {
		return -1, err
	}
	err = address_helper.CheckAddress(toAddress)
	if err != nil {
		return -1, err
	}

	var newBalance int

	err = d.Database.Transaction(func(tx *gorm.DB) error {
		var fromWallet *model.Wallet
		var toWallet *model.Wallet
		var err error

		// To avoid deadlocks, the wallets are queried in specific order.
		// Lexicographically smaller wallet is queried first. This guarantees
		// that no cycles of dependencies will occur.
		if fromAddress < toAddress {
			fromWallet, err = d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, fromAddress)
			if err != nil {
				return err
			}

			toWallet, err = d.getToWallet(ctx, tx, toAddress)
			if err != nil {
				return err
			}
		} else { // toAddress < fromAddress
			toWallet, err = d.getToWallet(ctx, tx, toAddress)
			if err != nil {
				return err
			}

			fromWallet, err = d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, fromAddress)
			if err != nil {
				return err
			}
		}

		if fromWallet.Tokens < amount {
			return errors.New("insufficient balance")
		}

		newFromWalletBalance := fromWallet.Tokens - amount
		newToWalletBalance := toWallet.Tokens + amount

		// Since both records are locked, there is no need to stick to the order any longer.
		err = d.WalletRepository.UpdateWalletTokensByAddress(ctx, tx, fromAddress, newFromWalletBalance)
		if err != nil {
			return err
		}

		err = d.WalletRepository.UpdateWalletTokensByAddress(ctx, tx, toAddress, newToWalletBalance)
		if err != nil {
			return err
		}

		newBalance = newFromWalletBalance
		return nil
	})
	if err != nil {
		return -1, err
	}

	return newBalance, nil
}

func (d *WalletService) getToWallet(ctx context.Context, tx *gorm.DB, toAddress string) (*model.Wallet, error) {
	toWallet, err := d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, toAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = d.WalletRepository.AddWallet(ctx, tx, &model.Wallet{
				Address: toAddress,
				Tokens:  0,
			})
			if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, err
			}

			toWallet, err = d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, toAddress)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return toWallet, err
}
