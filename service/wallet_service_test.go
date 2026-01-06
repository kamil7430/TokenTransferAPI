package service

import (
	"context"
	"testing"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/repository"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestWalletService(t *testing.T) {
	ctx := context.Background()
	dbname := "serviceTests"
	dbuser := "user"
	dbpassword := "password"

	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	dbURL, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := gorm.Open(gormpostgres.Open(dbURL), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Wallet{})
	require.NoError(t, err)

	d := WalletService{
		WalletRepository: &repository.DatabaseWalletRepository{},
		Database:         db,
	}

	t.Run("create wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")

		wallet, err := d.TryCreateWallet(ctx, "0x0000000000000000000000000000000000000001", 100)
		require.NoError(t, err)
		require.Equal(t, 100, wallet.Tokens)
		require.Equal(t, "0x0000000000000000000000000000000000000001", wallet.Address)

		wallet, err = d.TryCreateWallet(ctx, "0x0000000000000000000000000000000000000001", 150)
		require.Error(t, err)
	})

	t.Run("get wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)

		wallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		require.Equal(t, 100, wallet.Tokens)
		require.Equal(t, "0x0000000000000000000000000000000000000001", wallet.Address)
	})

	t.Run("transfer", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 200)

		amount, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 60)
		require.NoError(t, err)
		require.Equal(t, 40, amount)

		fromWallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000001", fromWallet.Address)
		require.Equal(t, 40, fromWallet.Tokens)

		toWallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000002", toWallet.Address)
		require.Equal(t, 260, toWallet.Tokens)
	})
}
