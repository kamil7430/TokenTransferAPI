package service

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
)

type WalletServicer interface {
	GetWallet(ctx context.Context, address string) (*model.Wallet, error)
	Transfer(ctx context.Context, fromAddress string, toAddress string, amount int) (int, error)
}
