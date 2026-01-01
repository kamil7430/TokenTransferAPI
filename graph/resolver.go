package graph

import "github.com/kamil7430/TokenTransferAPI/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	WalletService service.WalletServicer
}
