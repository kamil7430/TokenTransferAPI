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

	err := address_helper.CheckAddress(fromAddress)
	if err != nil {
		return -1, err
	}
	err = address_helper.CheckAddress(toAddress)
	if err != nil {
		return -1, err
	}

	// To avoid deadlocks, the wallets are queried in specific order.
	// Lexicographically smaller wallet is queried first. This guarantees
	// that no cycles of dependencies will occur.
	firstAddress, secondAddress := fromAddress, toAddress
	addressesFlipped := false
	if firstAddress > toAddress {
		firstAddress, secondAddress = secondAddress, firstAddress
		addressesFlipped = true
	}

	var newBalance int

	err = d.Database.Transaction(func(tx *gorm.DB) error {
		// Since it's not specified in the task, I assume that both wallets should exist on transfer.
		firstWallet, intErr := d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, firstAddress)
		if intErr != nil {
			return intErr // rollback on any error
		}

		secondWallet, intErr := d.WalletRepository.GetWalletByAddressForUpdate(ctx, tx, secondAddress)
		if intErr != nil {
			return intErr
		}

		// Since both records are locked, there is no need to stick to the order any longer.
		fromWallet, toWallet := firstWallet, secondWallet
		if addressesFlipped {
			fromWallet, toWallet = toWallet, fromWallet
		}

		if fromWallet.Tokens < amount {
			return errors.New("insufficient balance")
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

func (d *WalletService) TryCreateWallet(ctx context.Context, address string, tokens int) (*model.Wallet, error) {
	err := address_helper.CheckAddress(address)
	if err != nil {
		return nil, err
	}

	newWallet := &model.Wallet{
		Address: address,
		Tokens:  tokens,
	}

	err = d.WalletRepository.AddWallet(ctx, d.Database, newWallet)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("wallet with this address already exists")
		}
		return nil, err
	}

	return newWallet, nil
}
