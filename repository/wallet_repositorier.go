package repository

import (
	"context"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"gorm.io/gorm"
)

type WalletRepositorier interface { // Strange interface naming convention in Go
	GetWalletByAddress(ctx context.Context, tx *gorm.DB, address string) (*model.Wallet, error)
	GetWalletByAddressForUpdate(ctx context.Context, tx *gorm.DB, address string) (*model.Wallet, error)
	UpdateWalletTokensByAddress(ctx context.Context, tx *gorm.DB, address string, tokens string) error
}
