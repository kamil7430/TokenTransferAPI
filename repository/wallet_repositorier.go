package repository

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
)

type WalletRepositorier interface { // Strange interface naming convention in Go
	GetWalletByAddress(ctx context.Context, address string) (*model.Wallet, error)
	UpdateWalletByAddress(ctx context.Context, address string, wallet *model.Wallet) error
}
