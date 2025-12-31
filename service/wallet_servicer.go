package service

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
)

type WalletServicer interface {
	GetWalletByAddress(ctx context.Context, address string) (*model.Wallet, error)
	Transfer(ctx context.Context, fromAddress string, toAddress string, amount string) (string, error)
}
